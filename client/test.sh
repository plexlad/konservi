#!/usr/bin/env bash
set -e

echo "=== Running Frontend Tests ==="

# Install dependencies
pnpm install

# Run tests
pnpm test

# Run lint
pnpm lint

echo "=== Tests Complete ==="
