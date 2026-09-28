package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuntimeStateUpdate(t *testing.T) {
	latestState.Lock()
	latestState.value = nil
	latestState.Unlock()

	h := Handler()
	req := httptest.NewRequest(http.MethodGet, "/v0/update?event=after-adventure&character=tester&total_turns=123&ascension_turns=45&adventures=67&ascensions=8&breakfast=true&daily_sidecar=true&persistent_sidecar=false", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, body = %s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/v0/state", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("state status = %d", rec.Code)
	}
	var payload struct {
		Available bool `json:"available"`
		State     struct {
			Character          string `json:"character"`
			TotalTurns         int    `json:"total_turns"`
			BreakfastCompleted bool   `json:"breakfast_completed"`
			DailySidecar       bool   `json:"daily_sidecar"`
			PersistentSidecar  bool   `json:"persistent_sidecar"`
		} `json:"state"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if !payload.Available || payload.State.Character != "tester" || payload.State.TotalTurns != 123 || !payload.State.BreakfastCompleted || !payload.State.DailySidecar || payload.State.PersistentSidecar {
		t.Fatalf("unexpected state: %#v", payload)
	}
}
