#!/bin/bash
set -e

echo "=== Starting server with Streamable HTTP ==="
./bin/server -http -listen :9090 &
SERVER_PID=$!
sleep 1

echo "=== Listing tools via HTTP ==="
./bin/client -addr http://localhost:9090/mcp -tools-list

echo "=== Calling hello_world via HTTP ==="
./bin/client -addr http://localhost:9090/mcp -tool-call hello_world

echo "=== Stopping server ==="
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null || true

echo "=== All tests passed ==="
