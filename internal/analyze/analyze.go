package analyze

import (
	"sort"

	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
)

type Report struct {
	Backend          string         `json:"backend"`
	Variables        []Variable     `json:"variables"`
	Calls            []Call         `json:"calls"`
	Effects          []string       `json:"effects"`
	RawStatements    int            `json:"raw_statements"`
	UnsupportedLines []int          `json:"unsupported_lines,omitempty"`
	Summary          map[string]any `json:"summary"`
}

type Variable struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Call struct {
	Name   string `json:"name"`
	Effect string `json:"effect"`
	Count  int    `json:"count"`
}

func Program(p *ir.Program) Report {
	r := Report{Backend: p.Backend, Summary: map[string]any{}}
	counts := map[string]int{}
	effects := map[string]bool{}
	for _, s := range p.Statements {
		if s.Kind == "raw_statement" {
			r.RawStatements++
			r.UnsupportedLines = append(r.UnsupportedLines, s.Line)
		}
		if s.Var != nil {
			r.Variables = append(r.Variables, Variable{Name: s.Var.Name, Type: s.Var.Type.Name})
			if s.Var.Value != nil {
				walk(*s.Var.Value, counts)
			}
		}
		if s.Expr != nil {
			walk(*s.Expr, counts)
		}
	}
	for name, count := range counts {
		effect := effectOf(name)
		r.Calls = append(r.Calls, Call{Name: name, Count: count, Effect: effect})
		effects[effect] = true
	}
	sort.Slice(r.Calls, func(i, j int) bool { return r.Calls[i].Name < r.Calls[j].Name })
	sort.Slice(r.Variables, func(i, j int) bool { return r.Variables[i].Name < r.Variables[j].Name })
	for e := range effects {
		r.Effects = append(r.Effects, e)
	}
	sort.Strings(r.Effects)
	r.Summary["call_count"] = len(r.Calls)
	r.Summary["variable_count"] = len(r.Variables)
	r.Summary["fully_lowered"] = r.RawStatements == 0
	return r
}

func walk(e ir.Expr, counts map[string]int) {
	if e.Kind == "call" {
		counts[e.Name]++
	}
	for _, a := range e.Args {
		walk(a, counts)
	}
	if e.Left != nil {
		walk(*e.Left, counts)
	}
	if e.Right != nil {
		walk(*e.Right, counts)
	}
}

// This table is intentionally small and conservative. Unknown host calls are
// labeled unknown, never silently classified as read-only.
func effectOf(name string) string {
	switch name {
	case "my_hp", "my_maxhp", "my_mp", "my_meat", "item_amount", "have_item", "have_effect", "get_property", "available_amount", "print":
		return "read-or-output"
	case "adventure", "adv1", "use", "equip", "create", "buy", "sell", "cli_execute", "set_property", "visit_url":
		return "potential-mutation"
	default:
		return "unknown"
	}
}
