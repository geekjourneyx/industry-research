# researcher Makefile

VERSION ?= $(shell tr -d '[:space:]' < VERSION)
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: all build fast release clean test fmt vet install deps help

all: build

build:
	@echo "Building researcher..."
	@echo "Version: $(VERSION)"
	@go build -trimpath -ldflags="$(LDFLAGS)" -o researcher ./cmd/researcher
	@echo "Built: ./researcher"

fast:
	@go build -trimpath -ldflags="$(LDFLAGS)" -o researcher ./cmd/researcher

release:
	@echo "Building researcher release binaries..."
	@echo "Version: $(VERSION)"
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-linux-amd64 ./cmd/researcher
	@GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-linux-arm64 ./cmd/researcher
	@GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-darwin-amd64 ./cmd/researcher
	@GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-darwin-arm64 ./cmd/researcher
	@GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-windows-amd64.exe ./cmd/researcher
	@chmod +x bin/*-linux* bin/*-darwin* 2>/dev/null || true
	@ls -lh bin/

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
	@echo "  make release  - build release binaries into bin/"
	@echo "  make clean    - remove build outputs"
	@echo "  make test     - run tests"
	@echo "  make fmt      - format Go code"
	@echo "  make vet      - run go vet"
	@echo "  make install  - install to GOPATH/bin"
	@echo "  make deps     - download and tidy dependencies"
