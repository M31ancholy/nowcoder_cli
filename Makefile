.PHONY: build run clean test lint

BINARY_NAME=nowcoder_cli
VERSION?=v0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

build:
	go build -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" -o bin/${BINARY_NAME} ./cmd/cli

run:
	go run ./cmd/cli

clean:
	rm -rf bin/

test:
	go test -v ./...

lint:
	golangci-lint run

install:
	go install ./cmd/cli

revive:
	revive -config revive.toml ./...
