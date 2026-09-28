package parser

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
)

var varDeclRE = regexp.MustCompile(`^(void|boolean|int|float|string|buffer|item|effect|skill|location|monster|class|stat|familiar|slot|element|coinmaster)\s+([A-Za-z_][A-Za-z0-9_]*)\s*(?:=\s*(.+))?;$`)

func parseBootstrap(source []byte) (*ir.Program, error) {
	program := &ir.Program{Backend: "bootstrap"}
	s := bufio.NewScanner(strings.NewReader(string(source)))
	lineNo := 0
	for s.Scan() {
		lineNo++
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "import ") {
			program.Statements = append(program.Statements, ir.RawStatement(line, lineNo, "import lowering is not implemented in bootstrap mode"))
			continue
		}
		if m := varDeclRE.FindStringSubmatch(line); m != nil {
			stmt := ir.Statement{Kind: "variable_declaration", Line: lineNo, Var: &ir.VarDecl{Type: ir.TypeRef{Name: m[1]}, Name: m[2]}}
			if strings.TrimSpace(m[3]) != "" {
				expr, err := parseExpr(strings.TrimSpace(m[3]))
				if err != nil {
					program.Statements = append(program.Statements, ir.RawStatement(line, lineNo, "bootstrap expression parser: "+err.Error()))
					continue
				}
				stmt.Var.Value = expr
			}
			program.Statements = append(program.Statements, stmt)
			continue
		}
		if strings.HasSuffix(line, ";") {
			exprText := strings.TrimSpace(strings.TrimSuffix(line, ";"))
			expr, err := parseExpr(exprText)
			if err == nil {
				program.Statements = append(program.Statements, ir.Statement{Kind: "expression_statement", Expr: expr, Line: lineNo})
				continue
			}
		}
		program.Statements = append(program.Statements, ir.RawStatement(line, lineNo, "unsupported bootstrap syntax; preserved verbatim"))
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return program, nil
}
