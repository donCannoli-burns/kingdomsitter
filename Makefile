.PHONY: test smoke build grammar-install grammar-generate grammar-test ts-install ts-check clean

build:
	mkdir -p bin
	go build -o bin/kingdomsitter ./cmd/kingdomsitter

test:
	go test ./...

smoke: test build
	./bin/kingdomsitter chain
	./bin/kingdomsitter analyze examples/hello.ash
	./bin/kingdomsitter emit-ts examples/hello.ash > /tmp/kingdomsitter-hello.ts
	@echo "wrote /tmp/kingdomsitter-hello.ts"

grammar-install:
	cd tree-sitter-ash && npm install

grammar-generate:
	cd tree-sitter-ash && npm run generate

grammar-test:
	cd tree-sitter-ash && npm test

ts-install:
	cd ts-runtime && npm install

ts-check:
	cd ts-runtime && npm run check

clean:
	rm -rf bin dist tree-sitter-ash/src tree-sitter-ash/node_modules ts-runtime/node_modules
