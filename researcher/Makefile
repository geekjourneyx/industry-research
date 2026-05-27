# researcher Makefile

VERSION ?= $(shell tr -d '[:space:]' < VERSION)
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: all build fast release release-check clean test fmt vet install deps help

all: build

build:
	@echo "Building researcher..."
	@echo "Version: $(VERSION)"
	@go build -trimpath -ldflags="$(LDFLAGS)" -o researcher ./cmd/researcher
	@echo "Built: ./researcher"

fast:
	@go build -trimpath -ldflags="$(LDFLAGS)" -o researcher ./cmd/researcher

release:
	@../researcher/scripts/package-release.sh

release-check:
	@../researcher/scripts/check-release.sh

clean:
	@rm -f researcher
	@rm -rf bin
	@rm -f *.log

test:
	@go test -count=1 ./...

fmt:
	@go fmt ./...
	@gofmt -w .

vet:
	@go vet ./...

install:
	@go install ./cmd/researcher

deps:
	@go mod download
	@go mod tidy

help:
	@echo "researcher Makefile commands:"
	@echo "  make build    - build current platform binary"
	@echo "  make fast     - fast current platform build"
	@echo "  make release  - build release archives into ../dist/"
	@echo "  make release-check - validate docs, changelog, and tag before release"
	@echo "  make clean    - remove build outputs"
	@echo "  make test     - run tests"
	@echo "  make fmt      - format Go code"
	@echo "  make vet      - run go vet"
	@echo "  make install  - install to GOPATH/bin"
	@echo "  make deps     - download and tidy dependencies"
