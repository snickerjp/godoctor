#!/bin/bash
set -e

PORT=9091
CONTAINER_NAME=godoctor-test

echo "=== Starting godoctor container ==="
sudo docker run -d --rm --name $CONTAINER_NAME -p $PORT:8080 godoctor
sleep 2

echo "=== Listing tools ==="
./bin/client -addr http://localhost:$PORT/mcp -tools-list

echo "=== Calling read_docs for fmt.Println ==="
./bin/client -addr http://localhost:$PORT/mcp -tool-call read_docs fmt.Println

echo "=== Stopping container ==="
sudo docker stop $CONTAINER_NAME

echo "=== All tests passed ==="
