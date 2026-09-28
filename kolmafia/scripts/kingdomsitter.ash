// Kingdomsitter KoLmafia controller.
//
// Preferences:
//   _sidecar=true  -> daily mode. Custom underscore prefs clear at rollover/new ascension.
//   sidecar=true   -> persistent override; keep state sync enabled 24/7.
//
// Daily mode is deliberately quiet before breakfast. After breakfast completes,
// KoLmafia begins sending lightweight read-only state snapshots after adventures.
// The Go sidecar never receives execution authority from these preferences.

string KINGDOMSITTER_URL = "http://127.0.0.1:10423";

boolean kingdomsitter_pref_true(string name) {
    return get_property(name) == "true";
}

boolean kingdomsitter_daily_mode() {
    return kingdomsitter_pref_true("_sidecar");
}

boolean kingdomsitter_persistent_mode() {
    return kingdomsitter_pref_true("sidecar");
}

boolean kingdomsitter_breakfast_done() {
    return kingdomsitter_pref_true("breakfastCompleted");
}

boolean kingdomsitter_should_sync() {
    return kingdomsitter_persistent_mode() || (kingdomsitter_daily_mode() && kingdomsitter_breakfast_done());
}

boolean kingdomsitter_reachable() {
    buffer response = visit_url(KINGDOMSITTER_URL + "/health", false);
    return contains_text(response.to_string(), "\"ok\":true");
}

void kingdomsitter_mark_reachable(boolean reachable) {
    if (reachable) {
        set_property("_kingdomsitterUnavailable", "false");
    } else {
        set_property("_kingdomsitterUnavailable", "true");
    }
}

void kingdomsitter_reset_preferences() {
    set_property("_sidecar", "true");
    set_property("sidecar", "false");
    print("Kingdomsitter preferences reset: _sidecar=true, sidecar=false", "navy");
}

void kingdomsitter_print_preferences() {
    print("_sidecar=" + get_property("_sidecar") + "  (daily; clears at rollover/new ascension)", "navy");
    print("sidecar=" + get_property("sidecar") + "  (persistent override)", "navy");
    print("breakfastCompleted=" + get_property("breakfastCompleted"), "navy");
}

void kingdomsitter_help() {
    print("Kingdomsitter — Go + Tree-sitter ASH sidecar", "blue");
    print("", "black");
    print("Commands:", "black");
    print("  kingdomsitter help          Show this help", "black");
    print("  kingdomsitter chain         Show the runtime chain", "black");
    print("  kingdomsitter status        Show prefs + local daemon health", "black");
    print("  kingdomsitter prefs         Show the two user-facing prefs", "black");
    print("  kingdomsitter reset         _sidecar=true; sidecar=false", "black");
    print("  kingdomsitter daily-on      Enable daily mode for this KoL day", "black");
    print("  kingdomsitter daily-off     Disable daily mode for this KoL day", "black");
    print("  kingdomsitter always-on     Persistent 24/7 override", "black");
    print("  kingdomsitter always-off    Clear persistent override", "black");
    print("  kingdomsitter hooks-install Install login + after-adventure hooks if free", "black");
    print("  kingdomsitter hooks-remove  Remove only Kingdomsitter-owned hooks", "black");
    print("", "black");
    print("OS sidecar:", "black");
    print("  kingdomsitter start         Detached; terminal returns immediately", "black");
    print("  kingdomsitter stop", "black");
    print("  kingdomsitter conda-install Auto-start detached daemon on Conda activation", "black");
    print("", "black");
    print("v0 execution authority: false", "black");
}

void kingdomsitter_chain() {
    print("go.mod -> detached Go sidecar -> Tree-sitter ASH -> typed ASH IR -> TypeScript/agent -> Libram -> KoLmafia -> state sync -> Go sidecar", "blue");
}

void kingdomsitter_send_update(string event_name) {
    if (!kingdomsitter_should_sync()) {
        return;
    }

    if (kingdomsitter_pref_true("_kingdomsitterUnavailable")) {
        return;
    }

    string url = KINGDOMSITTER_URL + "/v0/update"
        + "?event=" + url_encode(event_name)
        + "&character=" + url_encode(my_name())
        + "&total_turns=" + to_string(total_turns_played())
        + "&ascension_turns=" + to_string(my_turncount())
        + "&adventures=" + to_string(my_adventures())
        + "&ascensions=" + to_string(my_ascensions())
        + "&breakfast=" + get_property("breakfastCompleted")
        + "&daily_sidecar=" + get_property("_sidecar")
        + "&persistent_sidecar=" + get_property("sidecar");

    buffer response = visit_url(url, false, true);
    if (!contains_text(response.to_string(), "\"accepted\":true")) {
        kingdomsitter_mark_reachable(false);
    }
}

void kingdomsitter_status() {
    kingdomsitter_print_preferences();

    if (kingdomsitter_reachable()) {
        kingdomsitter_mark_reachable(true);
        print("Kingdomsitter sidecar ON — local daemon reachable.", "navy");
        if (kingdomsitter_persistent_mode()) {
            print("Mode: persistent 24/7 override.", "navy");
        } else if (kingdomsitter_daily_mode() && !kingdomsitter_breakfast_done()) {
            print("Mode: daily gate armed; state sync waits for breakfast.", "navy");
        } else if (kingdomsitter_daily_mode()) {
            print("Mode: daily state sync active.", "navy");
        } else {
            print("Mode: daemon available, KoLmafia state sync disabled.", "navy");
        }
        return;
    }

    kingdomsitter_mark_reachable(false);
    print("Kingdomsitter sidecar OFF — local daemon is not reachable.", "navy");
}

void kingdomsitter_login() {
    if (kingdomsitter_persistent_mode()) {
        if (kingdomsitter_reachable()) {
            kingdomsitter_mark_reachable(true);
            print("Kingdomsitter sidecar ON — persistent override.", "navy");
            kingdomsitter_send_update("login");
        } else {
            kingdomsitter_mark_reachable(false);
            print("Kingdomsitter sidecar should be ON, but the daemon is not reachable.", "red");
            print("Try: kingdomsitter start", "red");
            print("If missing: go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest", "red");
        }
        return;
    }

    if (!kingdomsitter_daily_mode()) {
        print("Kingdomsitter sidecar OFF — daily gate is not enabled.", "navy");
        return;
    }

    if (!kingdomsitter_breakfast_done()) {
        print("Kingdomsitter sidecar OFF — daily mode is armed; waiting for breakfast.", "navy");
        return;
    }

    if (kingdomsitter_reachable()) {
        kingdomsitter_mark_reachable(true);
        print("Kingdomsitter sidecar ON — breakfast complete; daily state sync active.", "navy");
        kingdomsitter_send_update("login-after-breakfast");
        return;
    }

    kingdomsitter_mark_reachable(false);
    print("Kingdomsitter sidecar should now be ON, but the daemon is not reachable.", "red");
    print("Try: kingdomsitter start", "red");
    print("Or install/update it: go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest", "red");
}

void kingdomsitter_turn() {
    kingdomsitter_send_update("after-adventure");
}

void kingdomsitter_hooks_install() {
    string login_hook = get_property("loginScript");
    string turn_hook = get_property("afterAdventureScript");

    if (login_hook != "" && login_hook != "kingdomsitter login") {
        print("Not changing loginScript; it is already set to: " + login_hook, "red");
        print("Chain 'kingdomsitter login' from your existing login hook instead.", "red");
        return;
    }

    if (turn_hook != "" && turn_hook != "kingdomsitter turn") {
        print("Not changing afterAdventureScript; it is already set to: " + turn_hook, "red");
        print("Chain 'kingdomsitter turn' from your existing after-adventure hook instead.", "red");
        return;
    }

    set_property("loginScript", "kingdomsitter login");
    set_property("afterAdventureScript", "kingdomsitter turn");
    set_property("relayRunsAfterAdventureScript", "true");
    print("Kingdomsitter hooks installed.", "navy");
    print("loginScript=kingdomsitter login", "navy");
    print("afterAdventureScript=kingdomsitter turn", "navy");
    print("relayRunsAfterAdventureScript=true", "navy");
}

void kingdomsitter_hooks_remove() {
    if (get_property("loginScript") == "kingdomsitter login") {
        set_property("loginScript", "");
    }
    if (get_property("afterAdventureScript") == "kingdomsitter turn") {
        set_property("afterAdventureScript", "");
    }
    print("Removed Kingdomsitter-owned login/after-adventure hooks.", "navy");
    print("relayRunsAfterAdventureScript was left unchanged to avoid clobbering another script's preference.", "navy");
}

void main(string command) {
    string cmd = to_lower_case(command);

    if (cmd == "" || cmd == "help" || cmd == "--help" || cmd == "-h") {
        kingdomsitter_help();
        return;
    }

    if (cmd == "chain") {
        kingdomsitter_chain();
        return;
    }

    if (cmd == "status") {
        kingdomsitter_status();
        return;
    }

    if (cmd == "prefs") {
        kingdomsitter_print_preferences();
        return;
    }

    if (cmd == "reset" || cmd == "defaults") {
        kingdomsitter_reset_preferences();
        return;
    }

    if (cmd == "daily-on") {
        set_property("_sidecar", "true");
        kingdomsitter_mark_reachable(true);
        print("Kingdomsitter daily mode armed for this KoL day.", "navy");
        return;
    }

    if (cmd == "daily-off") {
        set_property("_sidecar", "false");
        print("Kingdomsitter daily mode OFF for this KoL day.", "navy");
        return;
    }

    if (cmd == "always-on") {
        set_property("sidecar", "true");
        kingdomsitter_mark_reachable(true);
        print("Kingdomsitter persistent sidecar override ON.", "navy");
        return;
    }

    if (cmd == "always-off") {
        set_property("sidecar", "false");
        print("Kingdomsitter persistent sidecar override OFF.", "navy");
        return;
    }

    if (cmd == "login") {
        kingdomsitter_login();
        return;
    }

    if (cmd == "turn") {
        kingdomsitter_turn();
        return;
    }

    if (cmd == "hooks-install") {
        kingdomsitter_hooks_install();
        return;
    }

    if (cmd == "hooks-remove") {
        kingdomsitter_hooks_remove();
        return;
    }

    print("Unknown Kingdomsitter command: " + command, "red");
    kingdomsitter_help();
}
