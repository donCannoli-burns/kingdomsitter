//go:build !treesitter

package parser

import "github.com/donCannoli-burns/kingdomsitter/internal/ir"

func Backend() string { return "bootstrap" }

func Parse(source []byte) (*ir.Program, error) {
	return parseBootstrap(source)
}
