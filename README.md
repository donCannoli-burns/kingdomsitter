# kingdomsitter

`kingdomsitter` is a Go-first ASH language sidecar for KoLmafia.

The intended runtime chain is:

```text
go.mod / kingdomsitter
        ↓
Go sidecar
        ↓
Tree-sitter ASH parser
        ↓
typed ASH IR
        ↓
TypeScript emitter + agent-readable analysis
        ↓
Libram / kolmafia TypeScript surface
        ↓
KoLmafia
        ↓
state/results back to the sidecar
```

The project deliberately separates **language understanding** from **game authority**:

- Go owns parsing orchestration, IR, analysis, agent envelopes, transpilation, testing, and process supervision.
- Tree-sitter owns ASH concrete syntax parsing.
- The typed IR is the stable boundary between source syntax and downstream tools.
- TypeScript is the interoperability layer.
- Libram provides higher-level KoL automation helpers.
- KoLmafia remains authoritative for real game state and real mutations.

## Status

This is a runnable bootstrap, not a claim of complete ASH compatibility.

The repository compiles immediately using a deliberately small bootstrap parser. The Tree-sitter grammar lives in `tree-sitter-ash/`; after generating `src/parser.c`, build with the `treesitter` tag to make Tree-sitter the syntax-validation front end:

```bash
cd tree-sitter-ash
npm install
npm run generate
npm test
cd ..
go test -tags treesitter ./...
go build -tags treesitter -o bin/kingdomsitter ./cmd/kingdomsitter
```

The bootstrap parser and the Tree-sitter front end both lower into the same Go IR. That lets the project grow the grammar without coupling every consumer to Tree-sitter node shapes.

## KoLmafia checkout

Kingdomsitter is packaged so KoLmafia can install its ASH-facing shim without projecting the Go/Tree-sitter development tree into the runtime directories. The root `manifest.json` points KoLmafia at `kolmafia/`.

From the KoLmafia gCLI:

```text
git checkout https://github.com/donCannoli-burns/kingdomsitter.git main
```

Then:

```text
kingdomsitter help
kingdomsitter chain
kingdomsitter status
```

The checkout does **not** turn the Go program into an in-process KoLmafia component. Install/start the sidecar separately:

```bash
go install github.com/donCannoli-burns/kingdomsitter/cmd/kingdomsitter@latest
kingdomsitter serve
```

See [`docs/KOLMAFIA-CHECKOUT.md`](docs/KOLMAFIA-CHECKOUT.md) for the exact install/update boundary.

## Quick start

```bash
go test ./...
go build -o kingdomsitter ./cmd/kingdomsitter

./kingdomsitter chain
./kingdomsitter serve
./kingdomsitter parse examples/hello.ash
./kingdomsitter analyze examples/hello.ash
./kingdomsitter emit-ts examples/hello.ash
```

Example input:

```ash
int hp = my_hp();
int teeth = item_amount($item[seal tooth]);
print("hp=" + hp + ", seal teeth=" + teeth);
```

Representative TypeScript output:

```ts
import { itemAmount, myHp, print } from "kolmafia";
import { $item } from "libram";

export function main(): void {
  const hp: number = myHp();
  const teeth: number = itemAmount($item`seal tooth`);
  print("hp=" + hp + ", seal teeth=" + teeth);
}
```

## Commands

```text
kingdomsitter chain
kingdomsitter parse <file.ash>
kingdomsitter analyze <file.ash>
kingdomsitter emit-ts <file.ash>
```

`serve` starts the analysis/transpilation sidecar on `127.0.0.1:10423` by default. That port belongs to **kingdomsitter**, not to stock KoLmafia. The v0 server intentionally exposes no execution endpoint.

HTTP surface:

```text
GET  /health
POST /v0/parse
POST /v0/analyze
POST /v0/agent
POST /v0/emit-ts
```

`parse` emits the stable Go IR as JSON. `analyze` emits an agent-oriented summary of symbols, host calls, inferred effects, and unsupported/raw statements. `emit-ts` creates a TypeScript target using the KoLmafia JS module and Libram typed constants.

## Parser strategy

The intended parser pipeline is:

```text
source.ash
   ↓
Tree-sitter CST
   ↓
ASH lowering
   ↓
kingdomsitter IR
```

During bootstrap, an intentionally conservative parser provides enough syntax to keep the Go module buildable before generated Tree-sitter C sources have been checked in. Unsupported syntax is retained as `raw_statement` rather than guessed at.

That is deliberate: **unknown ASH should degrade into visible unsupported IR, not silently change meaning.**

## Host boundary

The sidecar host API is intentionally transport-independent:

```go
type Host interface {
    Call(ctx context.Context, name string, args []Value) (Value, error)
}
```

A future KoLmafia adapter can use whichever integration surface is appropriate without changing the parser, analyzer, emitter, or agent protocol.

Potential adapters:

```text
MockHost          deterministic tests
RecordingHost     fixture generation / replay
KoLmafiaHost      live game-state bridge
StaticHost        analysis-only refusal of runtime calls
```

No stock KoLmafia TCP service or magic port is assumed by this repository.

## TypeScript runtime

`ts-runtime/` contains the downstream shape expected by generated code. It depends on `kolmafia` typings and Libram and defines an agent envelope rather than embedding an LLM provider.

This keeps:

```text
agent intelligence != execution authority
```

The agent consumes structured analysis; execution still travels through the explicit runtime/host layer.

## Repository map

```text
cmd/kingdomsitter/        CLI
internal/ir/              stable typed ASH IR
internal/parser/          bootstrap + optional Tree-sitter front end
internal/analyze/         agent-facing static analysis
internal/tsgen/           ASH IR → TypeScript emitter
internal/host/            host/runtime boundary
internal/agent/           structured agent envelope

tree-sitter-ash/          ASH grammar source + tests + Go binding shim
ts-runtime/               Libram/KoLmafia TS integration surface
examples/                 smoke-test ASH inputs
```

## Non-goals for v0

- pretending the grammar already covers every ASH production;
- reimplementing KoLmafia game state in Go;
- giving an LLM direct mutation authority;
- translating unknown constructs by guessing;
- treating Libram as the parser or semantic authority.

## Design invariant

The key architectural invariant is:

> Parse once into a stable typed ASH representation; let execution, analysis, migration, tests, and agents consume that representation independently.
