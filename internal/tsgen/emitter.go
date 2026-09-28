package tsgen

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
)

type Emitter struct {
	kolImports map[string]bool
	libImports map[string]bool
	warnings   []string
}

func Emit(p *ir.Program) string {
	e := &Emitter{kolImports: map[string]bool{}, libImports: map[string]bool{}}
	body := []string{}
	for _, s := range p.Statements {
		body = append(body, e.statement(s))
	}
	var out []string
	if len(e.kolImports) > 0 {
		out = append(out, fmt.Sprintf("import { %s } from \"kolmafia\";", strings.Join(sortedKeys(e.kolImports), ", ")))
	}
	if len(e.libImports) > 0 {
		out = append(out, fmt.Sprintf("import { %s } from \"libram\";", strings.Join(sortedKeys(e.libImports), ", ")))
	}
	if len(out) > 0 {
		out = append(out, "")
	}
	out = append(out, "export function main(): void {")
	for _, line := range body {
		for _, sub := range strings.Split(line, "\n") {
			out = append(out, "  "+sub)
		}
	}
	out = append(out, "}")
	if len(e.warnings) > 0 {
		out = append(out, "", "// kingdomsitter warnings:")
		for _, w := range e.warnings {
			out = append(out, "// - "+w)
		}
	}
	return strings.Join(out, "\n") + "\n"
}

func (e *Emitter) statement(s ir.Statement) string {
	switch s.Kind {
	case "variable_declaration":
		if s.Var == nil {
			return "// malformed variable declaration"
		}
		tsType := e.mapType(s.Var.Type.Name)
		if s.Var.Value == nil {
			return fmt.Sprintf("let %s: %s;", s.Var.Name, tsType)
		}
		return fmt.Sprintf("const %s: %s = %s;", s.Var.Name, tsType, e.expr(*s.Var.Value))
	case "expression_statement":
		if s.Expr == nil {
			return "// malformed expression statement"
		}
		return e.expr(*s.Expr) + ";"
	case "raw_statement":
		e.warnings = append(e.warnings, fmt.Sprintf("line %d was preserved rather than translated", s.Line))
		return "// TODO(kingdomsitter): unsupported ASH: " + strings.ReplaceAll(s.Raw, "*/", "* / ")
	default:
		return "// TODO(kingdomsitter): unknown IR statement " + s.Kind
	}
}

func (e *Emitter) expr(x ir.Expr) string {
	switch x.Kind {
	case "string", "number", "boolean":
		return x.Value
	case "identifier":
		return x.Name
	case "typed_constant":
		name := "$" + x.TypeName
		e.libImports[name] = true
		return fmt.Sprintf("%s`%s`", name, strings.ReplaceAll(x.Value, "`", "\\`"))
	case "call":
		name := snakeToCamel(x.Name)
		e.kolImports[name] = true
		args := make([]string, 0, len(x.Args))
		for _, a := range x.Args {
			args = append(args, e.expr(a))
		}
		return fmt.Sprintf("%s(%s)", name, strings.Join(args, ", "))
	case "binary":
		return fmt.Sprintf("%s %s %s", e.expr(*x.Left), x.Operator, e.expr(*x.Right))
	case "unary":
		return x.Operator + e.expr(*x.Right)
	default:
		return "undefined /* unsupported expression */"
	}
}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		r := []rune(parts[i])
		r[0] = unicode.ToUpper(r[0])
		parts[i] = string(r)
	}
	return strings.Join(parts, "")
}

func (e *Emitter) mapType(t string) string {
	switch t {
	case "int", "float":
		return "number"
	case "boolean":
		return "boolean"
	case "string", "buffer":
		return "string"
	case "item":
		e.kolImports["Item"] = true
		return "Item"
	case "effect":
		e.kolImports["Effect"] = true
		return "Effect"
	case "skill":
		e.kolImports["Skill"] = true
		return "Skill"
	case "location":
		e.kolImports["Location"] = true
		return "Location"
	case "monster":
		e.kolImports["Monster"] = true
		return "Monster"
	default:
		if t == "void" {
			return "void"
		}
		return "unknown"
	}
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
