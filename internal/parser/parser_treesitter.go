//go:build treesitter

package parser

import (
	"fmt"

	"github.com/donCannoli-burns/kingdomsitter/internal/ir"
	tree_sitter_ash "github.com/donCannoli-burns/kingdomsitter/tree-sitter-ash/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func Backend() string { return "tree-sitter+bootstrap-lowering" }

// Parse makes Tree-sitter the syntax authority, then lowers the accepted source
// into the stable kingdomsitter IR. The lowering pass is intentionally still
// conservative in v0 and can preserve unsupported constructs as raw statements.
func Parse(source []byte) (*ir.Program, error) {
	p := tree_sitter.NewParser()
	defer p.Close()
	if err := p.SetLanguage(tree_sitter.NewLanguage(tree_sitter_ash.Language())); err != nil {
		return nil, fmt.Errorf("set ASH language: %w", err)
	}
	tree := p.Parse(source, nil)
	if tree == nil {
		return nil, fmt.Errorf("tree-sitter returned no parse tree")
	}
	defer tree.Close()
	root := tree.RootNode()
	if root.HasError() {
		return nil, fmt.Errorf("ASH syntax error: %s", root.ToSexp())
	}
	program, err := parseBootstrap(source)
	if err != nil {
		return nil, err
	}
	program.Backend = Backend()
	return program, nil
}
