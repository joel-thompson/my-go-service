#!/bin/bash

set -e

case "$1" in
  setup)
    echo "Setting up development environment..."
    go mod tidy
    docker compose pull
    ;;
  
  start)
    echo "Starting API server..."
    if [ -f .env ]; then
      export $(cat .env | grep -v '^#' | xargs)
    fi
    go run ./cmd/server
    ;;
  
  migrate-up)
    echo "Running database migrations..."
    docker run --rm -v $(pwd)/migrations:/migrations --network host \
      migrate/migrate -path=/migrations/sql -database="postgres://postgres:postgres@localhost:5433/mygoservice?sslmode=disable" up
    ;;
  
  migrate-down)
    echo "Rolling back database migrations..."
    docker run --rm -v $(pwd)/migrations:/migrations --network host \
      migrate/migrate -path=/migrations/sql -database="postgres://postgres:postgres@localhost:5433/mygoservice?sslmode=disable" down
    ;;
  
  test)
    echo "Running tests..."
    go test ./...
    ;;
  
  build)
    echo "Building application..."
    go build -o bin/server ./cmd/server
    ;;
  
  build-cli)
    echo "Building CLI tool..."
    go build -o bin/mycli ./cmd/cli
    ;;
  
  lint)
    echo "Running linter..."
    go vet ./...
    go fmt ./...
    ;;
  
  *)
    echo "Usage: $0 {setup|start|migrate-up|migrate-down|test|build|lint}"
    exit 1
esac
