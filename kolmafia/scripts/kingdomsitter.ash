// Kingdomsitter KoLmafia checkout shim.
//
// This script intentionally has no game-mutation behavior. It only prints
// runtime information or performs a read-only HTTP GET against the local
// Kingdomsitter sidecar health endpoint.

string KINGDOMSITTER_URL = "http://127.0.0.1:10423";

void kingdomsitter_help() {
    print("Kingdomsitter — Go + Tree-sitter ASH sidecar", "blue");
    print("", "black");
    print("Commands:", "black");
    print("  kingdomsitter help    Show this help", "black");
    print("  kingdomsitter chain   Show the runtime chain", "black");
    print("  kingdomsitter status  GET the local sidecar /health endpoint", "black");
    print("", "black");
    print("Sidecar: " + KINGDOMSITTER_URL, "black");
    print("v0.1.0 execution authority: false", "black");
    print("Start the Go process separately with: kingdomsitter serve", "black");
}

void kingdomsitter_chain() {
    print("go.mod -> Go sidecar -> Tree-sitter ASH -> typed ASH IR -> TypeScript/agent -> Libram -> KoLmafia -> host adapter -> Go sidecar", "blue");
}

void kingdomsitter_status() {
    print("Checking " + KINGDOMSITTER_URL + "/health ...", "blue");
    buffer response = visit_url(KINGDOMSITTER_URL + "/health", false);
    print(response.to_string(), "black");
}

void main(string command) {
    string cmd = command.to_lower_case();

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

    print("Unknown Kingdomsitter command: " + command, "red");
    kingdomsitter_help();
}
