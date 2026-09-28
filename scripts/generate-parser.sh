#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../tree-sitter-ash"
if command -v tree-sitter >/dev/null 2>&1; then
  tree-sitter generate
else
  npx tree-sitter generate
fi
printf 'Generated tree-sitter-ash/src/parser.c\n'
printf 'Now run: go test -tags treesitter ./...\n'
