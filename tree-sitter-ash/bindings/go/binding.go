//go:build treesitter

package tree_sitter_ash

// #cgo CFLAGS: -std=c11 -fPIC
// #include "../../src/parser.c"
import "C"

import "unsafe"

// Language returns the generated Tree-sitter ASH language pointer.
// Run `npm run generate` in tree-sitter-ash before building with -tags treesitter.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_ash())
}
