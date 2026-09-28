package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/donCannoli-burns/kingdomsitter/internal/agent"
	"github.com/donCannoli-burns/kingdomsitter/internal/analyze"
	"github.com/donCannoli-burns/kingdomsitter/internal/parser"
	"github.com/donCannoli-burns/kingdomsitter/internal/tsgen"
)

const maxSourceBytes = 2 << 20

type runtimeSnapshot struct {
	UpdatedAt          string `json:"updated_at"`
	Event              string `json:"event"`
	Character          string `json:"character"`
	TotalTurns         int    `json:"total_turns"`
	AscensionTurns     int    `json:"ascension_turns"`
	Adventures         int    `json:"adventures"`
	Ascensions         int    `json:"ascensions"`
	BreakfastCompleted bool   `json:"breakfast_completed"`
	DailySidecar       bool   `json:"daily_sidecar"`
	PersistentSidecar  bool   `json:"persistent_sidecar"`
}

var latestState struct {
	sync.RWMutex
	value *runtimeSnapshot
}

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		latestState.RLock()
		seen := latestState.value != nil
		lastUpdate := ""
		if latestState.value != nil {
			lastUpdate = latestState.value.UpdatedAt
		}
		latestState.RUnlock()
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":                  true,
			"name":                "kingdomsitter",
			"parser_backend":      parser.Backend(),
			"execution_authority": false,
			"state_seen":          seen,
			"last_state_update":   lastUpdate,
		})
	})
	mux.HandleFunc("GET /v0/state", func(w http.ResponseWriter, _ *http.Request) {
		latestState.RLock()
		defer latestState.RUnlock()
		if latestState.value == nil {
			writeJSON(w, http.StatusOK, map[string]any{"available": false})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"available":           true,
			"execution_authority": false,
			"state":               latestState.value,
		})
	})
	mux.HandleFunc("GET /v0/update", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		snapshot := &runtimeSnapshot{
			UpdatedAt:          time.Now().UTC().Format(time.RFC3339),
			Event:              q.Get("event"),
			Character:          q.Get("character"),
			TotalTurns:         queryInt(q.Get("total_turns")),
			AscensionTurns:     queryInt(q.Get("ascension_turns")),
			Adventures:         queryInt(q.Get("adventures")),
			Ascensions:         queryInt(q.Get("ascensions")),
			BreakfastCompleted: queryBool(q.Get("breakfast")),
			DailySidecar:       queryBool(q.Get("daily_sidecar")),
			PersistentSidecar:  queryBool(q.Get("persistent_sidecar")),
		}
		latestState.Lock()
		latestState.value = snapshot
		latestState.Unlock()
		writeJSON(w, http.StatusOK, map[string]any{
			"accepted":            true,
			"execution_authority": false,
			"state":               snapshot,
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

func queryInt(value string) int {
	n, _ := strconv.Atoi(value)
	return n
}

func queryBool(value string) bool {
	b, _ := strconv.ParseBool(value)
	return b
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
