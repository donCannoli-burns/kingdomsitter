package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenKind int

const (
	tEOF tokenKind = iota
	tIdent
	tNumber
	tString
	tTypedConstant
	tLParen
	tRParen
	tComma
	tOperator
)

type token struct {
	kind tokenKind
	text string
}

func lexExpr(src string) ([]token, error) {
	var out []token
	for i := 0; i < len(src); {
		c := src[i]
		if unicode.IsSpace(rune(c)) {
			i++
			continue
		}
		if c == '$' {
			start := i
			i++
			for i < len(src) && (unicode.IsLetter(rune(src[i])) || src[i] == '_') {
				i++
			}
			if i >= len(src) || src[i] != '[' {
				return nil, fmt.Errorf("invalid typed constant near %q", src[start:])
			}
			depth := 0
			for i < len(src) {
				if src[i] == '[' {
					depth++
				} else if src[i] == ']' {
					depth--
					if depth == 0 {
						i++
						break
					}
				}
				i++
			}
			if depth != 0 {
				return nil, fmt.Errorf("unclosed typed constant")
			}
			out = append(out, token{kind: tTypedConstant, text: src[start:i]})
			continue
		}
		if unicode.IsLetter(rune(c)) || c == '_' {
			start := i
			i++
			for i < len(src) && (unicode.IsLetter(rune(src[i])) || unicode.IsDigit(rune(src[i])) || src[i] == '_') {
				i++
			}
			out = append(out, token{kind: tIdent, text: src[start:i]})
			continue
		}
		if unicode.IsDigit(rune(c)) {
			start := i
			i++
			for i < len(src) && (unicode.IsDigit(rune(src[i])) || src[i] == '.') {
				i++
			}
			out = append(out, token{kind: tNumber, text: src[start:i]})
			continue
		}
		if c == '"' {
			start := i
			i++
			escaped := false
			for i < len(src) {
				if escaped {
					escaped = false
					i++
					continue
				}
				if src[i] == '\\' {
					escaped = true
					i++
					continue
				}
				if src[i] == '"' {
					i++
					break
				}
				i++
			}
			if i > len(src) || src[i-1] != '"' {
				return nil, fmt.Errorf("unclosed string literal")
			}
			out = append(out, token{kind: tString, text: src[start:i]})
			continue
		}
		switch c {
		case '(':
			out = append(out, token{kind: tLParen, text: "("})
			i++
		case ')':
			out = append(out, token{kind: tRParen, text: ")"})
			i++
		case ',':
			out = append(out, token{kind: tComma, text: ","})
			i++
		default:
			ops := []string{"==", "!=", ">=", "<=", "&&", "||", "+", "-", "*", "/", ">", "<", "!"}
			matched := ""
			for _, op := range ops {
				if strings.HasPrefix(src[i:], op) {
					matched = op
					break
				}
			}
			if matched == "" {
				return nil, fmt.Errorf("unexpected character %q", c)
			}
			out = append(out, token{kind: tOperator, text: matched})
			i += len(matched)
		}
	}
	out = append(out, token{kind: tEOF})
	return out, nil
}
