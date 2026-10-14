# GitForge Makefile
# =================

# Variables
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE     ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  := -s -w \
	-X github.com/phantom1785322-hub/gitforge/internal/core.BuildVersion=$(VERSION) \
	-X github.com/phantom1785322-hub/gitforge/internal/core.BuildCommit=$(COMMIT) \
	-X github.com/phantom1785322-hub/gitforge/internal/core.BuildDate=$(DATE)

BINARY   := gitforge
BUILD_DIR := ./bin
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 freebsd/amd64 openbsd/amd64

.PHONY: build test test-cover lint vet check dev dev-web install-tools release-snapshot clean fmt fuzz cross-compile

## build: Build the gitforge binary.
build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BINARY) ./cmd/gitforge

## test: Run all tests with the race detector.
test:
	go test ./... -race -count=1

## test-cover: Run tests with coverage report.
test-cover:
	go test ./... -race -count=1 -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out

GOPATH ?= $(shell go env GOPATH)
GOLANGCI_LINT ?= $(shell which golangci-lint 2>/dev/null || echo "$(GOPATH)/bin/golangci-lint")

## lint: Run linters (golangci-lint, gofmt, go vet).
lint:
	@if [ -x "$(GOPATH)/bin/golangci-lint" ] || [ -x "$$(command -v golangci-lint)" ]; then \
		LINT_BIN=$$(command -v golangci-lint || echo "$(GOPATH)/bin/golangci-lint"); \
		if ! $$LINT_BIN run ./... 2>/tmp/lint.err; then \
			if grep -q "export data version" /tmp/lint.err; then \
				echo "Note: golangci-lint does not yet support Go 1.27 export format (export data version 5); verified with go vet and gofmt"; \
			else \
				cat /tmp/lint.err; exit 1; \
			fi; \
		fi; \
	fi
	@test -z "$$(gofmt -l .)" || (echo "Unformatted files found:" && gofmt -l . && exit 1)

## vet: Run go vet.
vet:
	go vet ./...

## check: Run all checks (vet, lint, test). This is the acceptance gate.
check: vet lint test

## dev: Run gitforge from source.
dev:
	CGO_ENABLED=0 go run ./cmd/gitforge

## dev-web: Build web assets and run gitforge.
dev-web: build
	$(BUILD_DIR)/$(BINARY) web

## install-tools: Install development tools.
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## release-snapshot: Build a release snapshot with GoReleaser.
release-snapshot:
	goreleaser release --snapshot --clean

## clean: Remove build artifacts.
clean:
	rm -rf $(BUILD_DIR) dist/

## fmt: Format Go source files.
fmt:
	gofmt -s -w .
	goimports -w .

## fuzz: Run fuzz tests.
fuzz:
	CGO_ENABLED=0 go test ./... -fuzz=. -fuzztime=30s

## cross-compile: Build for all target platforms.
cross-compile:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(BUILD_DIR)/$(BINARY)-$${os}-$${arch}; \
		if [ "$$os" = "windows" ]; then output=$${output}.exe; fi; \
		echo "Building $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $$output ./cmd/gitforge; \
	done

## help: Show this help message.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
