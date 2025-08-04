# My Go Service

A production-ready Go backend service built with Gin, PostgreSQL, and industry best practices.

## Quick Start

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- PostgreSQL (via Docker)

### Setup
```bash
# Clone and setup
git clone <repository>
cd my-go-service

# Start database
docker compose up -d

# Setup dependencies and run migrations  
./do setup
./do migrate-up

# Start API server
./do start
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Health check |
| GET    | `/hello`  | Hello world |  
| POST   | `/items`  | Create item |
| GET    | `/items`  | List items with pagination |
| GET    | `/items/:id` | Get single item by ID |
| PUT    | `/items/:id` | Update existing item |
| DELETE | `/items/:id` | Delete item |

### CLI Testing Tool

Build and use the CLI for easy API testing:

```bash
# Build CLI
./do build-cli

# Test endpoints
./bin/mycli health
./bin/mycli hello

# Items CRUD operations
./bin/mycli items create --name "Test Item" --description "My item"
./bin/mycli items list
./bin/mycli items get --id <item-id>
./bin/mycli items update --id <item-id> --name "Updated Name"
./bin/mycli items delete --id <item-id>

# Pagination
./bin/mycli items list --limit 5 --offset 10

# Custom server URL
./bin/mycli --url http://localhost:8080 health

# JSON output
./bin/mycli --format json items list

# Verbose mode
./bin/mycli -v items create --name "Debug Item"
```

### Development Commands

```bash
./do setup          # Install dependencies
./do start          # Start API server  
./do build          # Build API server
./do build-cli      # Build CLI tool
./do migrate-up     # Run migrations
./do migrate-down   # Rollback migrations
./do test           # Run tests
./do lint           # Run linter
```

### Configuration

Copy `.env.example` to `.env` and customize:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5433/mygoservice?sslmode=disable
SERVER_ADDR=:8080
LOG_LEVEL=info
```

## Project Structure

```
my-go-service/
├── cmd/
│   ├── server/         # API server
│   └── cli/           # CLI testing tool  
├── api/server/        # HTTP handlers & routing
├── storage/           # Database layer
├── constants/         # Shared constants
├── migrations/sql/    # Database migrations
└── .env              # Local configuration
```

## Learning Objectives

This project demonstrates:
- **Go Fundamentals**: Packages, interfaces, error handling
- **Web APIs**: REST endpoints with Gin framework
- **Database**: PostgreSQL with migrations and SQLX
- **CLI Tools**: Cobra framework with subcommands
- **Configuration**: Environment-based config management
- **Development**: Local setup with Docker & task automation

## Next Steps

1. Add more endpoints (`GET /items` with pagination)
2. Add comprehensive tests
3. Add external service clients
4. Consider deployment options