package tsgen

import (
	"strings"
	"testing"

	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
)

func TestEmit(t *testing.T) {
	p := &ir.Program{Statements: []ir.Statement{
		{Kind: "variable_declaration", Var: &ir.VarDecl{Type: ir.TypeRef{Name: "int"}, Name: "hp", Value: &ir.Expr{Kind: "call", Name: "my_hp"}}},
		{Kind: "expression_statement", Expr: &ir.Expr{Kind: "call", Name: "print", Args: []ir.Expr{{Kind: "string", Value: "\"hi\""}}}},
	}}
	got := Emit(p)
	if !strings.Contains(got, "myHp") || !strings.Contains(got, "export function main") {
		t.Fatal(got)
	}
}
