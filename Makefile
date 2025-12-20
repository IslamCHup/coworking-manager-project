# Makefile for Go project
MAIN ?= ./cmd/api


run:
	go run $(MAIN)

air:
	@command -v air >/dev/null 2>&1 || { $(MAKE) install-air; }
	air

install-air:
	go install github.com/cosmtrek/air@latest

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) $(MAIN)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { $(MAKE) install-linter; }
	golangci-lint run

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

deps:
	go mod tidy

clean:
	rm -rf bin