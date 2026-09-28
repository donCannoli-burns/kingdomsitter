package analyze

import (
	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
	"testing"
)

func TestUnknownIsNotReadOnly(t *testing.T) {
	p := &ir.Program{Statements: []ir.Statement{{Kind: "expression_statement", Expr: &ir.Expr{Kind: "call", Name: "mystery_call"}}}}
	r := Program(p)
	if len(r.Calls) != 1 || r.Calls[0].Effect != "unknown" {
		t.Fatalf("unexpected classification: %#v", r.Calls)
	}
}
