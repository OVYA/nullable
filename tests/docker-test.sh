#!/bin/bash
# Convenience script to build and run the Docker PostgreSQL test

set -e

echo "Building Docker image..."
docker build -t nullable-postgres-test -f tests/Dockerfile .

echo ""
echo "Running PostgreSQL integration test..."
docker run --rm nullable-postgres-test

echo ""
echo "✓ Test completed!"
