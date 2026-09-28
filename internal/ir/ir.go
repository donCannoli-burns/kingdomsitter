package ir

// Program is the stable language-neutral representation consumed by analysis,
// transpilation, tests, and agent tooling.
type Program struct {
	Backend    string      `json:"backend"`
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Kind     string   `json:"kind"`
	Var      *VarDecl `json:"var,omitempty"`
	Expr     *Expr    `json:"expr,omitempty"`
	Raw      string   `json:"raw,omitempty"`
	Line     int      `json:"line,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type VarDecl struct {
	Type  TypeRef `json:"type"`
	Name  string  `json:"name"`
	Value *Expr   `json:"value,omitempty"`
}

type TypeRef struct {
	Name string `json:"name"`
}

type Expr struct {
	Kind     string         `json:"kind"`
	Name     string         `json:"name,omitempty"`
	Value    string         `json:"value,omitempty"`
	TypeName string         `json:"type_name,omitempty"`
	Args     []Expr         `json:"args,omitempty"`
	Left     *Expr          `json:"left,omitempty"`
	Right    *Expr          `json:"right,omitempty"`
	Operator string         `json:"operator,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
}

func RawStatement(text string, line int, warning string) Statement {
	s := Statement{Kind: "raw_statement", Raw: text, Line: line}
	if warning != "" {
		s.Warnings = []string{warning}
	}
	return s
}
