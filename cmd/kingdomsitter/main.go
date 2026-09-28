package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/donCannoli-burns/kingdomsitter/internal/agent"
	"github.com/donCannoli-burns/kingdomsitter/internal/analyze"
	"github.com/donCannoli-burns/kingdomsitter/internal/daemon"
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
		addr := optionalAddr()
		fmt.Fprintf(os.Stderr, "kingdomsitter sidecar listening on http://%s (analysis/transpile only)\n", addr)
		must(http.ListenAndServe(addr, server.Handler()))
	case "start":
		addr := optionalAddr()
		message, err := daemon.Start(addr)
		must(err)
		fmt.Println(message)
	case "stop":
		if len(os.Args) != 2 {
			usage(2)
		}
		message, err := daemon.Stop()
		must(err)
		fmt.Println(message)
	case "status":
		addr := optionalAddr()
		message, err := daemon.Status(addr)
		must(err)
		fmt.Println(message)
	case "conda-install":
		if len(os.Args) != 2 {
			usage(2)
		}
		path, err := daemon.InstallCondaHook()
		must(err)
		fmt.Printf("installed Conda activation hook: %s\n", path)
	case "conda-remove":
		if len(os.Args) != 2 {
			usage(2)
		}
		path, err := daemon.RemoveCondaHook()
		must(err)
		fmt.Printf("removed Conda activation hook: %s\n", path)
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

func optionalAddr() string {
	if len(os.Args) == 2 {
		return daemon.DefaultAddr
	}
	if len(os.Args) == 3 {
		return os.Args[2]
	}
	usage(2)
	return ""
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
  kingdomsitter start [127.0.0.1:10423]      # detached; terminal returns
  kingdomsitter stop
  kingdomsitter status [127.0.0.1:10423]
  kingdomsitter serve [127.0.0.1:10423]      # foreground/debug
  kingdomsitter conda-install
  kingdomsitter conda-remove
  kingdomsitter parse <file.ash>
  kingdomsitter analyze <file.ash>
  kingdomsitter agent <file.ash>
  kingdomsitter emit-ts <file.ash>`)
	os.Exit(code)
}
