#!/usr/bin/env bash
set -e

echo "=== Starting Konservi Dev Environment ==="

# Check dependencies
command -v go >/dev/null 2>&1 || { echo "Go not found"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "Node not found"; exit 1; }
command -v pnpm >/dev/null 2>&1 || { echo "PNPM not found"; exit 1; }

# Start backend
echo "Starting backend with air..."
cd server && air &
BACKEND_PID=$!

# Start frontend
echo "Starting frontend with vite..."
cd ../frontend && pnpm dev &
FRONTEND_PID=$!

# Cleanup on exit
trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null" EXIT

# Wait for processes
wait
