export GOWORK ?= off

.PHONY: all test test-full test-race test-coverage lint build clean setup-hooks format vulncheck

all: lint test build

setup-hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/*
	@echo "Charites git hooks installed to .githooks"

# Fast test: lightweight iteration
test:
	@if [ -f "go.mod" ]; then \
		go test -v ./...; \
	else \
		echo "Notice: go.mod not initialized yet. Run Phase 0 setup first."; \
	fi

# Full test with Go Race Detector
test-full:
	@if [ -f "go.mod" ]; then \
		go test -race -v ./...; \
	else \
		echo "Notice: go.mod not initialized yet. Run Phase 0 setup first."; \
	fi

test-race: test-full

# Test Coverage
COVER_PKGS ?= github.com/will2469/charites/internal/...,github.com/will2469/charites/cmd/...
test-coverage:
	go test -race -coverpkg=$(COVER_PKGS) -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -func=coverage.txt | tail -n 1
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint:
	@if [ -f "go.mod" ]; then \
		echo "Running go vet..."; \
		go vet ./...; \
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

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)

build:
	go build -ldflags="$(LDFLAGS)" -v -o bin/charites ./cmd/charites

format:
	@if [ -x .githooks/format-all.sh ]; then \
		.githooks/format-all.sh; \
	else \
		gofmt -w $$(find . -name '*.go' -not -path './vendor/*' 2>/dev/null); \
	fi

clean:
	rm -rf bin/ dist/ coverage.txt coverage.html
