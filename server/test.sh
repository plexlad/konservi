#!/usr/bin/env bash
set -e

echo "=== Running Backend Tests ==="

go test -v -cover -race ./...

golangci-lint run

echo "=== Tests Complete ==="
