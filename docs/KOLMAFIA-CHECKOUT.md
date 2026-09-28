# KoLmafia checkout

Kingdomsitter is a Go sidecar with a small KoLmafia-facing shim.

## Install into KoLmafia

Run in the KoLmafia gCLI:

```text
git checkout https://github.com/donCannoli-burns/kingdomsitter.git main
```

The repository uses `manifest.json` with `root_directory: kolmafia`, so KoLmafia installs only the contents of `kolmafia/` into the user-facing KoLmafia directories. The Go, Tree-sitter, TypeScript, tests, and development files remain part of the checked-out Git working copy without being projected into `scripts/`, `data/`, or `relay/`.

After checkout, the gCLI command is:

```text
kingdomsitter help
kingdomsitter chain
kingdomsitter status
```

`kingdomsitter status` performs a read-only GET to `http://127.0.0.1:10423/health`.

## Install the Go sidecar

The KoLmafia checkout installs the ASH shim; the sidecar is a separate operating-system process. With Go installed:

```bash
go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest
kingdomsitter serve
```

Or build from a normal source checkout:

```bash
go build -o bin/kingdomsitter ./cmd/kingdomsitter
./bin/kingdomsitter serve
```

The default listener is loopback-only:

```text
127.0.0.1:10423
```

## Update

In KoLmafia:

```text
git update
```

For a `go install` installation, rebuild/update the binary separately:

```bash
go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest
```

## Boundary

The v0.1.0 sidecar exposes parsing, analysis, agent-envelope generation, and TypeScript emission. It intentionally exposes no game-execution endpoint and reports `execution_authority: false`.
