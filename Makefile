.DEFAULT_GOAL := build

.PHONY:fmt vet staticcheck build clean test

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

staticcheck: vet
	staticcheck ./...

build: staticcheck
	go build -o ./build/automata-compiler ./cmd/automata-compiler

clean:
	rm -rf build/*

test:
	go test -v -cover ./...