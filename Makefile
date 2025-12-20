# Makefile for Go project (PowerShell)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

MAIN ?= ./cmd/api

run:
    go run $(MAIN)

air:
    if (-not (Get-Command air -ErrorAction SilentlyContinue)) { $(MAKE) install-air }; air

install-air:
    go install github.com/cosmtrek/air@latest

build:
    New-Item -ItemType Directory -Force -Path bin | Out-Null
    go build -o bin/$(BINARY) $(MAIN)

test:
    go test ./...

fmt:
    go fmt ./...

vet:
    go vet ./...

lint:
    if (-not (Get-Command golangci-lint -ErrorAction SilentlyContinue)) { $(MAKE) install-linter }; golangci-lint run

install-linter:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

deps:
    go mod tidy

clean:
    if (Test-Path bin) { Remove-Item -Recurse -Force bin }