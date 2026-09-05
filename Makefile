export GOWORK ?= off

.PHONY: all test test-unit test-full test-race test-coverage lint build clean cross-compile setup-hooks format vulncheck wiki

# Deterministic build-before-test ordering (ARCH-00-BUILD-ORDER / TEST-00-BUILD-ORDER)
all: build test lint

setup-hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/*
	@echo "Charites git hooks installed to .githooks"

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X github.com/will2469/charites/internal/cli.Version=$(VERSION) -X github.com/will2469/charites/internal/cli.Commit=$(COMMIT) -X github.com/will2469/charites/internal/cli.BuildDate=$(BUILD_DATE)

# Native Pure Go compilation with Zero CGO (SPEC-00-BUILD-001)
build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/charites ./cmd/charites

# Concurrency verification test: runs Go Race Detector via ThreadSanitizer harness
test: build
	@if [ -f "go.mod" ]; then \
		CGO_ENABLED=1 go test -v -race ./...; \
	else \
		echo "Notice: go.mod not initialized yet. Run Phase 0 setup first."; \
	fi

# Pure Go Unit Test: runs in any environment with CGO_ENABLED=0 (zero C compiler requirement)
test-unit: build
	@if [ -f "go.mod" ]; then \
		go test -v ./...; \
	else \
		echo "Notice: go.mod not initialized yet. Run Phase 0 setup first."; \
	fi

# Full test with Go Race Detector
test-full: build
	@if [ -f "go.mod" ]; then \
		CGO_ENABLED=1 go test -race -v ./...; \
	else \
		echo "Notice: go.mod not initialized yet. Run Phase 0 setup first."; \
	fi

test-race: test-full

# Test Coverage
COVER_PKGS ?= github.com/will2469/charites/internal/...,github.com/will2469/charites/cmd/...
test-coverage: build
	CGO_ENABLED=1 go test -race -coverpkg=$(COVER_PKGS) -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -func=coverage.txt | tail -n 1
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Cross-platform compilation verification for 4 official release targets (SPEC-00-BUILD-002 / TEST-00-BUILD-002)
cross-compile:
	@echo "Verifying cross-compilation targets (SPEC-00-BUILD-002)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites

lint:
	@if [ -f "go.mod" ]; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	fi
	@echo "Checking gofmt..."
	@unformatted=$$(find . -name '*.go' -not -path './vendor/*' -exec gofmt -l {} + 2>/dev/null); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted files found: $$unformatted"; \
		exit 1; \
	fi
	@echo "Code hygiene clean!"

vulncheck:
	govulncheck ./...

format:
	@if [ -x .githooks/format-all.sh ]; then \
		.githooks/format-all.sh; \
	else \
		gofmt -w $$(find . -name '*.go' -not -path './vendor/*' 2>/dev/null); \
	fi

wiki:
	@echo "Regenerating wiki documentation from rules..."
	@go test -run TestGenerator_RegenerateWiki ./internal/wiki/... > /dev/null
	@echo "Wiki documentation generated successfully at wiki/"

clean:
	rm -rf bin/ dist/ coverage.txt coverage.html
