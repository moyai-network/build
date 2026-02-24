#!/usr/bin/env bash
set -euo pipefail

if ! command -v golangci-lint >/dev/null 2>&1; then
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
fi

if [ -f .golangci.yml ]; then
  golangci-lint run --config .golangci.yml
else
  golangci-lint run ./...
fi
