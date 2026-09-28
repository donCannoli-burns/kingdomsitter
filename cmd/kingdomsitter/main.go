package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/donCannoli-burns/kingdomsitter/internal/agent"
	"github.com/donCannoli-burns/kingdomsitter/internal/analyze"
	"github.com/donCannoli-burns/kingdomsitter/internal/parser"
	"github.com/donCannoli-burns/kingdomsitter/internal/server"
	"github.com/donCannoli-burns/kingdomsitter/internal/tsgen"
)

func main() {
	if len(os.Args) < 2 {
		usage(2)
	}
	switch os.Args[1] {
	case "chain":
		fmt.Printf("go.mod -> Go sidecar -> %s -> ASH IR -> TypeScript/agent -> Libram -> KoLmafia -> host adapter -> Go sidecar\n", parser.Backend())
	case "serve":
		addr := "127.0.0.1:10423"
		if len(os.Args) == 3 {
			addr = os.Args[2]
		} else if len(os.Args) > 3 {
			usage(2)
		}
		fmt.Fprintf(os.Stderr, "kingdomsitter sidecar listening on http://%s (analysis/transpile only)\n", addr)
		must(http.ListenAndServe(addr, server.Handler()))
	case "parse", "analyze", "agent", "emit-ts":
		if len(os.Args) != 3 {
			usage(2)
		}
		src, err := os.ReadFile(os.Args[2])
		must(err)
		program, err := parser.Parse(src)
		must(err)
		switch os.Args[1] {
		case "parse":
			writeJSON(program)
		case "analyze":
			writeJSON(analyze.Program(program))
		case "agent":
			writeJSON(agent.New("review-and-transform-ash", analyze.Program(program)))
		case "emit-ts":
			fmt.Print(tsgen.Emit(program))
		}
	default:
		usage(2)
	}
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	must(enc.Encode(v))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "kingdomsitter:", err)
		os.Exit(1)
	}
}

func usage(code int) {
	fmt.Fprintln(os.Stderr, `kingdomsitter - Go/Tree-sitter ASH sidecar

usage:
  kingdomsitter chain
  kingdomsitter serve [127.0.0.1:10423]
  kingdomsitter parse <file.ash>
  kingdomsitter analyze <file.ash>
  kingdomsitter agent <file.ash>
  kingdomsitter emit-ts <file.ash>`)
	os.Exit(code)
}
