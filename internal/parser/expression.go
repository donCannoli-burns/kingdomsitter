package parser

import (
	"fmt"
	"strings"

	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
)

type exprParser struct {
	tokens []token
	pos    int
}

func parseExpr(src string) (*ir.Expr, error) {
	tokens, err := lexExpr(src)
	if err != nil {
		return nil, err
	}
	p := &exprParser{tokens: tokens}
	expr, err := p.parseBinary(0)
	if err != nil {
		return nil, err
	}
	if p.peek().kind != tEOF {
		return nil, fmt.Errorf("unexpected trailing token %q", p.peek().text)
	}
	return &expr, nil
}

var precedence = map[string]int{
	"||": 1,
	"&&": 2,
	"==": 3, "!=": 3,
	">": 4, "<": 4, ">=": 4, "<=": 4,
	"+": 5, "-": 5,
	"*": 6, "/": 6,
}

func (p *exprParser) parseBinary(minPrec int) (ir.Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return ir.Expr{}, err
	}
	for {
		tok := p.peek()
		if tok.kind != tOperator {
			break
		}
		prec, ok := precedence[tok.text]
		if !ok || prec < minPrec {
			break
		}
		p.next()
		right, err := p.parseBinary(prec + 1)
		if err != nil {
			return ir.Expr{}, err
		}
		lcopy, rcopy := left, right
		left = ir.Expr{Kind: "binary", Operator: tok.text, Left: &lcopy, Right: &rcopy}
	}
	return left, nil
}

func (p *exprParser) parsePrimary() (ir.Expr, error) {
	tok := p.next()
	switch tok.kind {
	case tString:
		return ir.Expr{Kind: "string", Value: tok.text}, nil
	case tNumber:
		return ir.Expr{Kind: "number", Value: tok.text}, nil
	case tTypedConstant:
		open := strings.IndexByte(tok.text, '[')
		return ir.Expr{Kind: "typed_constant", TypeName: strings.TrimPrefix(tok.text[:open], "$"), Value: tok.text[open+1 : len(tok.text)-1]}, nil
	case tIdent:
		if p.peek().kind == tLParen {
			p.next()
			args := []ir.Expr{}
			if p.peek().kind != tRParen {
				for {
					arg, err := p.parseBinary(0)
					if err != nil {
						return ir.Expr{}, err
					}
					args = append(args, arg)
					if p.peek().kind != tComma {
						break
					}
					p.next()
				}
			}
			if p.next().kind != tRParen {
				return ir.Expr{}, fmt.Errorf("expected )")
			}
			return ir.Expr{Kind: "call", Name: tok.text, Args: args}, nil
		}
		switch tok.text {
		case "true", "false":
			return ir.Expr{Kind: "boolean", Value: tok.text}, nil
		default:
			return ir.Expr{Kind: "identifier", Name: tok.text}, nil
		}
	case tLParen:
		expr, err := p.parseBinary(0)
		if err != nil {
			return ir.Expr{}, err
		}
		if p.next().kind != tRParen {
			return ir.Expr{}, fmt.Errorf("expected )")
		}
		return expr, nil
	case tOperator:
		if tok.text == "!" || tok.text == "-" {
			right, err := p.parsePrimary()
			if err != nil {
				return ir.Expr{}, err
			}
			return ir.Expr{Kind: "unary", Operator: tok.text, Right: &right}, nil
		}
	}
	return ir.Expr{}, fmt.Errorf("unexpected token %q", tok.text)
}

func (p *exprParser) peek() token { return p.tokens[p.pos] }
func (p *exprParser) next() token {
	t := p.tokens[p.pos]
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	return t
}
