# Architecture Guide

A comprehensive guide to understanding the structure and architecture of this Go service.

## Overview

This Go service follows a **layered architecture** pattern with clear separation of concerns. It implements a RESTful API with PostgreSQL persistence, includes a CLI testing tool, and uses modern Go practices for configuration, logging, and database interactions.

## Architectural Patterns

- **Layered Architecture**: Clear separation between presentation (API), business logic, and data layers
- **Dependency Injection**: Dependencies are injected through constructors and interfaces
- **Command Pattern**: CLI commands are structured using the command pattern via Cobra
- **Repository Pattern**: Data access is abstracted through the storage layer
- **Configuration as Code**: Environment-based configuration with sensible defaults

## Project Structure

```
my-go-service/
├── cmd/                    # Application entry points
│   ├── server/            # HTTP API server
│   └── cli/               # CLI testing tool
├── api/server/            # HTTP layer (routes, handlers, middleware)
├── storage/               # Data access layer
├── constants/             # Shared application constants
├── migrations/sql/        # Database schema migrations
├── clients/               # External service clients (empty, for future use)
├── _plans/                # Project planning and enhancement documents
├── bin/                   # Compiled binaries
├── do                     # Build and development script
├── docker-compose.yml     # Local PostgreSQL setup
├── go.mod                 # Go module definition
└── README.md             # Quick start guide
```

## Core Components

### 1. Entry Points (`cmd/`)

#### Server Entry Point (`cmd/server/`)
- **`main.go`**: HTTP server entry point with graceful shutdown
  - Initializes application dependencies
  - Sets up HTTP server with Gin router
  - Handles OS signals for graceful shutdown (SIGTERM, SIGINT)
  - Implements 30-second shutdown timeout

- **`setup/setup.go`**: Application bootstrap and dependency injection
  - **Config struct**: Environment-based configuration
    - `SERVER_ADDR`: HTTP server bind address (default: `:8080`)
    - `DATABASE_URL`: PostgreSQL connection string (required)
    - `LOG_LEVEL`: Logging level (default: `info`)
  - **App struct**: Dependency container holding logger, database, and config
  - Handles database connection setup with connection pooling
  - Configures structured JSON logging with configurable levels

#### CLI Entry Point (`cmd/cli/`)
- **`main.go`**: Simple CLI entry point that delegates to Cobra commands
- **`commands/`**: CLI command implementations
  - **`root.go`**: Base command with global flags (`--url`, `--format`, `--verbose`)
  - **`health.go`**: Health check command (`mycli health`)
  - **`hello.go`**: Hello world command (`mycli hello`)
  - **`items.go`**: Complete CRUD operations for items
    - `create`: Create new items with validation
    - `list`: List items with pagination support
    - `get`: Retrieve single item by ID
    - `update`: Update existing items
    - `delete`: Delete items by ID
  - Consistent error handling across all commands
  - Support for JSON and pretty-print output formats
  - Connection error handling with user-friendly messages

### 2. HTTP Layer (`api/server/`)

#### API Server (`api.go`)
- **API struct**: Holds logger and storage dependencies
- **SetupRoutes()**: Configures Gin router with middleware and routes
  - Uses Gin release mode for production
  - Custom logging middleware for structured request logging
  - Recovery middleware for panic handling
- **Route Definitions**:
  - `GET /health`: Health check endpoint
  - `GET /hello`: Simple hello world endpoint
  - `POST /items`: Create new item
  - `GET /items`: List items with pagination
  - `GET /items/:id`: Get item by ID
  - `PUT /items/:id`: Update item
  - `DELETE /items/:id`: Delete item

#### Request Handlers (`handlers.go`)
- Implements all HTTP handlers following the `handleVerbNoun` naming pattern
- **Request Flow**:
  1. Parse and validate request parameters/body
  2. Call appropriate storage layer method
  3. Handle errors with appropriate HTTP status codes
  4. Return structured JSON responses
- **Error Handling**:
  - 400 Bad Request for validation errors
  - 404 Not Found for missing resources
  - 500 Internal Server Error for system errors
- **Features**:
  - UUID-based resource identification
  - Pagination support for list endpoints
  - Partial updates for PUT operations

### 3. Storage Layer (`storage/`)

#### Storage Interface (`store.go`)
- **Store struct**: Database operations wrapper around `*sqlx.DB`
- **Methods**: Complete CRUD operations with context support
  - `CreateItem()`: Insert new items with RETURNING clause
  - `ListItems()`: Paginated listing with total count
  - `GetItem()`: Single item retrieval by UUID
  - `UpdateItem()`: Partial updates using COALESCE
  - `DeleteItem()`: Soft delete returning deleted item
- **Features**:
  - Context-aware operations for cancellation/timeout
  - Automatic pagination defaults and limits
  - SQL injection prevention via parameterized queries
  - Proper error handling and type conversion

#### Data Models and Queries (`sql.go`)
- **Item struct**: Core data model with JSON and DB tags
  - UUID primary key with automatic generation
  - Nullable description field
  - Automatic timestamp management
- **Request/Response DTOs**:
  - `CreateItemRequest`: Input validation with required fields
  - `UpdateItemRequest`: Optional fields for partial updates
  - `ListItemsRequest`: Pagination parameters
  - `ListItemsResponse`: Paginated response with metadata
- **SQL Queries**: Raw SQL queries as constants
  - Parameterized queries for security
  - RETURNING clauses for atomic operations
  - COALESCE for partial updates
  - Proper indexing considerations

### 4. Configuration (`constants/`)

#### Shared Constants (`constants.go`)
- **HTTP Headers**: Content type definitions
- **Response Messages**: Standardized API responses
- **Status Codes**: Application-specific status constants
- Centralized location for magic strings and values

### 5. Database Layer (`migrations/`)

#### SQL Migrations (`migrations/sql/`)
- **Migration Files**: Versioned database schema changes
  - `000001_create_items_table.up.sql`: Initial items table creation
  - `000001_create_items_table.down.sql`: Rollback script
- **Schema Design**:
  - UUID primary keys for distributed systems
  - Timestamp columns with timezone support
  - Appropriate constraints and defaults
  - PostgreSQL-specific features (gen_random_uuid())

### 6. Development Tools

#### Build Script (`do`)
- **Bash script** providing consistent development commands
- **Commands**:
  - `setup`: Install dependencies and pull Docker images
  - `start`: Start API server with environment loading
  - `build`: Compile API server binary
  - `build-cli`: Compile CLI tool binary
  - `migrate-up/down`: Database migration management
  - `test`: Run Go tests
  - `lint`: Code formatting and linting
- **Environment Handling**: Automatic `.env` file loading
- **Docker Integration**: Uses migrate/migrate image for migrations

#### Docker Compose (`docker-compose.yml`)
- **PostgreSQL Service**: 
  - PostgreSQL 15 image
  - Port mapping (5433:5432) to avoid conflicts
  - Persistent volume for data
  - Health checks for startup coordination
  - Consistent database credentials

## Data Flow

### HTTP Request Flow
1. **Client Request** → Gin Router
2. **Middleware** → Logging + Recovery
3. **Handler** → Request validation and parsing
4. **Storage Layer** → Database operation
5. **Response** → JSON serialization and HTTP response

### CLI Command Flow
1. **CLI Input** → Cobra command parsing
2. **HTTP Client** → API server request
3. **Response Processing** → JSON parsing or error handling
4. **Output Formatting** → Pretty print or JSON output

### Database Operations
1. **Context Creation** → Request-scoped context
2. **Query Execution** → Parameterized SQL via sqlx
3. **Result Mapping** → Struct scanning with tags
4. **Error Handling** → Typed errors with context

## Key Design Decisions

### Technology Choices
- **Gin Framework**: High-performance HTTP router with middleware support
- **PostgreSQL**: ACID compliance, JSON support, and UUID generation
- **sqlx**: Raw SQL with struct scanning for performance and control
- **Cobra**: Feature-rich CLI framework with consistent UX
- **Structured Logging**: JSON logging for production observability
- **Environment Config**: 12-factor app configuration principles

### Architectural Benefits
- **Testability**: Clean interfaces and dependency injection
- **Maintainability**: Clear separation of concerns and consistent patterns
- **Observability**: Structured logging and error handling
- **Development Experience**: CLI tool for easy API testing
- **Production Ready**: Graceful shutdown, health checks, and error handling

### Trade-offs
- **Raw SQL vs ORM**: Chosen for performance and explicit control
- **Monolithic Structure**: Simple deployment vs microservice complexity
- **Environment Config**: Explicit configuration vs convention-based defaults

## Extension Points

### Adding New Endpoints
1. Add route in `api/server/api.go`
2. Implement handler in `api/server/handlers.go`
3. Add storage method in `storage/store.go`
4. Create CLI command in `cmd/cli/commands/`
5. Update documentation

### Adding External Services
- Create client interfaces in `clients/` directory
- Inject clients through `setup.App` struct
- Mock clients for testing

### Database Changes
- Create new migration files with sequential numbers
- Update data models in `storage/sql.go`
- Add corresponding CRUD operations

This architecture provides a solid foundation for a production-ready Go service while maintaining simplicity and clear patterns for future development.