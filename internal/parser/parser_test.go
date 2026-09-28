package parser

import "testing"

func TestBootstrapParse(t *testing.T) {
	src := []byte("int hp = my_hp();\nint n = item_amount($item[seal tooth]);\nprint(\"x=\" + n);\n")
	p, err := parseBootstrap(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Statements) != 3 {
		t.Fatalf("got %d statements", len(p.Statements))
	}
	if p.Statements[0].Var == nil || p.Statements[0].Var.Name != "hp" {
		t.Fatalf("bad first statement: %#v", p.Statements[0])
	}
	if p.Statements[1].Var == nil || p.Statements[1].Var.Value == nil || p.Statements[1].Var.Value.Args[0].Kind != "typed_constant" {
		t.Fatalf("typed constant did not parse")
	}
}
