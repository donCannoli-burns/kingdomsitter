package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/donCannoli-burns/kingdomsitter/internal/agent"
	"github.com/donCannoli-burns/kingdomsitter/internal/analyze"
	"github.com/donCannoli-burns/kingdomsitter/internal/parser"
	"github.com/donCannoli-burns/kingdomsitter/internal/tsgen"
)

const maxSourceBytes = 2 << 20

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":                  true,
			"name":                "kingdomsitter",
			"parser_backend":      parser.Backend(),
			"execution_authority": false,
		})
	})
	mux.HandleFunc("POST /v0/parse", sourceEndpoint(func(src []byte) (any, error) {
		return parser.Parse(src)
	}))
	mux.HandleFunc("POST /v0/analyze", sourceEndpoint(func(src []byte) (any, error) {
		p, err := parser.Parse(src)
		if err != nil {
			return nil, err
		}
		return analyze.Program(p), nil
	}))
	mux.HandleFunc("POST /v0/agent", sourceEndpoint(func(src []byte) (any, error) {
		p, err := parser.Parse(src)
		if err != nil {
			return nil, err
		}
		return agent.New("review-and-transform-ash", analyze.Program(p)), nil
	}))
	mux.HandleFunc("POST /v0/emit-ts", func(w http.ResponseWriter, r *http.Request) {
		src, err := readSource(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p, err := parser.Parse(src)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set("Content-Type", "text/typescript; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, tsgen.Emit(p))
	})
	return mux
}

func sourceEndpoint(fn func([]byte) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		src, err := readSource(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		v, err := fn(src)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		writeJSON(w, http.StatusOK, v)
	}
}

func readSource(r io.Reader) ([]byte, error) {
	limited := io.LimitReader(r, maxSourceBytes+1)
	b, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(b) > maxSourceBytes {
		return nil, fmt.Errorf("source exceeds %d bytes", maxSourceBytes)
	}
	return b, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
