# Go Service Blueprint

## Overview

This document outlines the plan for creating a Go backend service that mirrors the architecture and patterns of the `llmops-service` project. The goal is to learn Go by building a service with professional-grade structure and tooling.

## Project Goals

- Learn Go by building a real backend service
- Mirror the architecture patterns from `llmops-service`
- Start simple with basic endpoints and build up
- Use industry-standard tools and practices
- Keep everything runnable locally without external dependencies

## Architecture Blueprint

### Core Principles
1. **Structured Packages**: Clear separation of concerns across packages
2. **Dependency Injection**: Dependencies created at startup and injected where needed
3. **Configuration Management**: Environment-based configuration using structs
4. **Local Development**: Docker Compose for local dependencies
5. **Task Automation**: Shell script for common development tasks

### Directory Structure

```
my-go-service/
├── cmd/
│   └── server/
│       ├── main.go              # Application entry point
│       └── setup/
│           └── setup.go         # Dependency injection and app initialization
├── api/
│   └── server/
│       ├── api.go              # Router setup and API struct
│       └── handlers.go         # HTTP handlers for endpoints
├── storage/
│   ├── store.go                # Store struct and high-level methods
│   └── sql.go                  # Database models and raw SQL queries
├── constants/
│   └── constants.go            # Shared constants across packages
├── clients/                    # Future external service clients
│   └── (empty for now, will add openai/ later)
├── migrations/
│   └── sql/
│       ├── 000001_create_items_table.up.sql
│       └── 000001_create_items_table.down.sql
├── docker-compose.yml          # Local development environment
├── do                          # Task runner script
├── go.mod                      # Go module definition
└── README.md                   # Project documentation
```

## Technology Stack

### Core Framework
- **Web Framework**: `gin-gonic/gin` - High-performance HTTP web framework
- **Database Driver**: `github.com/jackc/pgx/v5` - PostgreSQL driver
- **Database Utilities**: `github.com/jmoiron/sqlx` - Enhanced database operations
- **Configuration**: `github.com/sethvargo/go-envconfig` - Environment variable parsing
- **Logging**: `log/slog` - Go's built-in structured logging
- **Database Migrations**: `golang-migrate/migrate` - Database schema versioning

### Development Tools
- **Database**: PostgreSQL 15 (via Docker)
- **Migration Runner**: golang-migrate (via Docker)
- **Task Runner**: Custom `./do` shell script

## API Endpoints (Initial)

### Health Check
- **Endpoint**: `GET /health`
- **Purpose**: Kubernetes/Docker health checks
- **Response**: `{"status": "healthy"}`

### Hello World
- **Endpoint**: `GET /hello`
- **Purpose**: Simple test endpoint
- **Response**: `{"message": "Hello, World!"}`

### Items Management
- **Endpoint**: `POST /items`
- **Purpose**: Create new items in database
- **Request**: `{"name": "string", "description": "string"}`
- **Response**: `{"id": "uuid", "name": "string", "description": "string", "created_at": "timestamp"}`

- **Endpoint**: `GET /items` (future)
- **Purpose**: List items with pagination
- **Response**: `{"items": [...], "total": number}`

## Database Schema

### Items Table
```sql
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Configuration Structure

```go
type Config struct {
    ServerAddr   string `env:"SERVER_ADDR,default=:8080"`
    DatabaseURL  string `env:"DATABASE_URL,required"`
    LogLevel     string `env:"LOG_LEVEL,default=info"`
}
```

Environment variables:
- `SERVER_ADDR`: HTTP server bind address (default: `:8080`)
- `DATABASE_URL`: PostgreSQL connection string
- `LOG_LEVEL`: Logging level (default: `info`)

## Development Workflow

### Setup Commands
```bash
./do setup          # Install development tools
./do migrate-up     # Run database migrations
./do start          # Start the API server
```

### Development Commands
```bash
./do test           # Run tests
./do lint           # Run linter
./do build          # Build binaries
./do migrate-down   # Rollback migrations
```

### Local Environment
```bash
docker compose up -d    # Start PostgreSQL
./do migrate-up         # Apply migrations
./do start              # Start API server
```

## Implementation Phases

### Phase 1: Foundation
1. Create project structure and `go.mod`
2. Set up `docker-compose.yml` with PostgreSQL
3. Create initial database migration
4. Implement basic HTTP server with Gin
5. Add configuration loading
6. Create `./do` script

### Phase 2: Core API
1. Implement health check endpoint
2. Implement hello world endpoint
3. Create storage layer (Store struct)
4. Implement items database model
5. Add POST /items endpoint

### Phase 3: Enhancement
1. Add GET /items endpoint with pagination
2. Add proper error handling and logging
3. Add tests for all components
4. Improve documentation

### Phase 4: Future Extensions
1. Add more complex business logic
2. Add external service clients (OpenAI, etc.)
3. Add more sophisticated endpoints
4. Consider deployment options

## Key Learning Objectives

### Go Language Concepts
- Package organization and visibility
- Struct embedding and composition
- Interface usage for abstraction
- Error handling patterns
- Context usage for request lifecycle

### Backend Development Patterns
- Dependency injection
- Repository pattern (storage layer)
- Configuration management
- Database migrations
- HTTP middleware
- Structured logging

### Development Practices
- Local development environment setup
- Database schema evolution
- Testing strategies
- Build and deployment automation
- Documentation practices

## Success Criteria

1. **Functional API**: All planned endpoints work correctly
2. **Professional Structure**: Code is well-organized and maintainable
3. **Local Development**: Easy setup for new developers
4. **Database Integration**: Proper data persistence and migrations
5. **Testing**: Comprehensive test coverage
6. **Documentation**: Clear setup and usage instructions

## Next Steps

1. Review and approve this blueprint
2. Create the project directory structure
3. Initialize Go module and dependencies
4. Set up local development environment
5. Begin Phase 1 implementation

---

*This blueprint serves as the foundation for building a production-ready Go service while learning industry best practices.*