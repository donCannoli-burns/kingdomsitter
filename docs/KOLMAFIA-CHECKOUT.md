# KoLmafia checkout

Kingdomsitter is a Go sidecar with a small KoLmafia-facing controller.

## Install into KoLmafia

Run in the KoLmafia gCLI:

```text
git checkout https://github.com/donCannoli-burns/kingdomsitter.git main
```

The repository uses `manifest.json` with `root_directory: kolmafia`, so KoLmafia installs only the contents of `kolmafia/` into the user-facing KoLmafia directories.

## Install the Go sidecar

With Go installed:

```bash
go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest
kingdomsitter start
```

`kingdomsitter start` is detached: it starts the local process, writes logs under the user's cache directory, and immediately returns the terminal. `kingdomsitter serve` still exists as a foreground/debug mode.

Useful OS commands:

```bash
kingdomsitter start
kingdomsitter status
kingdomsitter stop
```

The default listener is loopback-only:

```text
127.0.0.1:10423
```

## Conda activation trigger

If Kingdomsitter was installed while a Conda environment is active:

```bash
kingdomsitter conda-install
```

This writes an `activate.d/kingdomsitter.sh` hook into the current `$CONDA_PREFIX`. Every activation runs `kingdomsitter start`; start is idempotent and detached, so the terminal is not occupied.

Remove it with:

```bash
kingdomsitter conda-remove
```

The Conda hook only supervises process availability. It does **not** change KoLmafia preferences. KoLmafia owns the daily/persistent sync policy below.

## KoLmafia preferences

Kingdomsitter uses two user-facing custom properties:

| Preference | Meaning |
|---|---|
| `_sidecar` | Daily mode. A custom underscore preference clears automatically at rollover or new ascension. |
| `sidecar` | Persistent 24/7 override. |

Requested/default reset state:

```text
_sidecar=true
sidecar=false
```

Set it from gCLI with:

```text
kingdomsitter reset
```

or run the tiny reset tool directly:

```text
kingdomsitter_reset
```

Other controls:

```text
kingdomsitter daily-on
kingdomsitter daily-off
kingdomsitter always-on
kingdomsitter always-off
kingdomsitter prefs
kingdomsitter status
```

### Daily state machine

```text
sidecar=true
    -> persistent override
    -> state sync ON regardless of breakfast/daily reset

sidecar=false + _sidecar=true + breakfastCompleted=false
    -> daemon may exist in background
    -> KoLmafia state sync OFF
    -> navy login message

sidecar=false + _sidecar=true + breakfastCompleted=true
    -> daily state sync ON
    -> if daemon missing, red login diagnostic with install/start command

sidecar=false + _sidecar=false
    -> state sync OFF
```

The distinction is intentional: the detached Go daemon can remain cheaply available while KoLmafia decides whether live character state is being sent to it.

## Login + after-adventure hooks

Install the hooks with:

```text
kingdomsitter hooks-install
```

Kingdomsitter will only claim the hook slots if they are empty or already point to Kingdomsitter. It refuses to overwrite another script.

The installed values are:

```text
loginScript = kingdomsitter login
afterAdventureScript = kingdomsitter turn
relayRunsAfterAdventureScript = true
```

The after-adventure hook is KoLmafia's documented lifecycle surface for running a script after adventures. It is the closest safe approximation to "update every turn"; it is **not** advertised as a universal hook for every possible turn-consuming action.

Remove only Kingdomsitter-owned hooks with:

```text
kingdomsitter hooks-remove
```

## State sync

When sync is enabled, KoLmafia sends a lightweight loopback-only snapshot to:

```text
GET /v0/update
```

The snapshot currently contains character name, total turns, current-ascension turns, adventures remaining, ascension count, breakfast state, and the two sidecar preference values.

Read the latest snapshot locally:

```text
GET /v0/state
```

No execution command is accepted by this endpoint and the sidecar continues to report:

```text
execution_authority=false
```

## Update

In KoLmafia:

```text
git update
```

Update the Go binary separately:

```bash
go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest
kingdomsitter stop
kingdomsitter start
```
