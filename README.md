# MapBot Shared - Shared Go Utilities

[![Go Version](https://img.shields.io/github/go-mod/go-version/pixime/mapbot-shared)](https://go.dev/)
[![License](https://img.shields.io/github/license/pixime/mapbot-shared)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pixime/mapbot-shared)](https://goreportcard.com/report/github.com/pixime/mapbot-shared)

Shared Go utilities for the MapBot project ecosystem. This module provides common functionality used across multiple MapBot services including database management, configuration, logging, and testing utilities.

## 📦 Packages

### `database`

PostgreSQL database manager with connection pooling, health checks, and automatic migrations support.

```go
import "github.com/pixime/mapbot-shared/database"

cfg := config.NewPostgresDatabase("localhost", 5432, "mydb", "user", "pass")
dm, err := database.NewDatabaseManager(cfg,
    database.WithMigrations("./migrations"),
)
defer dm.Close()

// Use standard database/sql
db := dm.GetDB()

// Or use pgxpool for advanced features
pool := dm.GetPool()
```

**Features:**

- Connection pooling with configurable limits
- Health checks and statistics
- Automatic schema migrations via golang-migrate
- Support for both `database/sql` and `pgxpool`
- SSL/TLS support

### `config`

Configuration structures and utilities for PostgreSQL connections.

```go
import "github.com/pixime/mapbot-shared/config"

cfg := config.NewPostgresDatabase("localhost", 5432, "mydb", "user", "pass")
// Returns sensible defaults for production use

// Get connection string
connStr := cfg.ConnectionString()

// Get public connection string (password masked)
publicStr := cfg.PublicConnectionString()
```

### `logger`

Structured logging with `slog` and environment-based configuration.

```go
import "github.com/pixime/mapbot-shared/logger"

// Auto-configured via environment variables:
// LOG_LEVEL=debug|info|warn|error (default: info)
// LOG_FORMAT=json|text (default: text)

log := logger.GetLogger()
log.Info("Application started", "version", "1.0.0")
```

### `testutils`

Testing utilities for integration tests with testcontainers.

```go
import "github.com/pixime/mapbot-shared/testutils"

func TestMyDatabase(t *testing.T) {
    // Automatically spins up PostGIS container, runs migrations, and cleans up
    dm, cleanup := testutils.SetupTestDatabaseManager(t)
    defer cleanup()

    // Your tests here
    db := dm.GetDB()
    // ...
}

// With custom options
dm, cleanup := testutils.SetupTestDatabaseManager(t, testutils.TestDatabaseManagerOptions{
    MigrationsPath: "./custom/migrations",
    SkipMigrations: false,
    MigrationsSchema: "test_schema",
    MigrationsTable: "my_migrations",
})

// Environment setup for tests
testutils.SetupTestEnv(t)  // Loads .env.test

// With required variables
testutils.SetupTestEnvWithRequiredVarsOrSkipTest(t, "API_KEY", "DATABASE_URL")

// Get test data file path
path := testutils.GetTestDataFilePath("sample.json")
```

## 🚀 Installation

```bash
go get github.com/pixime/mapbot-shared@latest
```

## 📋 Requirements

- Go 1.23 or higher
- PostgreSQL 12+ (for database features)
- Docker (for running tests with testcontainers)

## �️ Development Commands

```bash
# Show all available commands
make help

# Build all packages
make build

# Run tests
make test

# Run tests with coverage report
make test-coverage

# View HTML coverage report
make coverage

# Run linter (golangci-lint)
make lint

# Install dependencies
make deps

# Clean generated files
make clean

# Run all checks and build
make all
```

## �🔧 Usage in Projects

### Standard Import

```go
import (
    "github.com/pixime/mapbot-shared/database"
    "github.com/pixime/mapbot-shared/config"
    "github.com/pixime/mapbot-shared/logger"
    "github.com/pixime/mapbot-shared/testutils"
)
```

### With Go Workspace (Local Development)

If you're working on multiple MapBot projects and want to modify `mapbot-shared` locally:

```bash
# In your workspace root
go work init
go work use ./mapbot-shared
go work use ./mapbot-ai
go work use ./french-admin-etl
go work sync
```

Now changes to `mapbot-shared` are immediately visible in other projects!

## 🧪 Running Tests

Using Makefile (recommended):

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View detailed HTML coverage report
make coverage
```

Or directly with Go:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Verbose output
go test -v ./...
```

**Note:** Database tests require Docker to be running for testcontainers.

## 📖 Examples

### Complete Database Setup

```go
package main

import (
    "context"
    "log"

    "github.com/pixime/mapbot-shared/config"
    "github.com/pixime/mapbot-shared/database"
)

func main() {
    // Configure database
    cfg := config.NewPostgresDatabase(
        "localhost",
        5432,
        "myapp",
        "user",
        "password",
    )
    cfg.MaxOpenConns = 50
    cfg.SSLMode = "require"

    // Initialize with automatic migrations
    dm, err := database.NewDatabaseManager(cfg,
        database.WithMigrations("./migrations"),
    )
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer dm.Close()

    // Health check
    ctx := context.Background()
    if err := dm.Health(ctx); err != nil {
        log.Fatalf("Database unhealthy: %v", err)
    }

    // Use the database
    db := dm.GetDB()
    rows, err := db.Query("SELECT id, name FROM users")
    // ...
}
```

### Testing with TestContainers

```go
package mypackage_test

import (
    "testing"

    "github.com/pixime/mapbot-shared/testutils"
)

func TestUserRepository(t *testing.T) {
    // Setup test database with migrations
    dm, cleanup := testutils.SetupTestDatabaseManager(t)
    defer cleanup()

    // Create repository
    repo := NewUserRepository(dm.GetDB())

    // Test your code
    err := repo.CreateUser(&User{Name: "Alice"})
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    users, err := repo.ListUsers()
    if err != nil {
        t.Fatalf("Failed to list users: %v", err)
    }

    if len(users) != 1 {
        t.Errorf("Expected 1 user, got %d", len(users))
    }
}
```

## 📝 Migration Files

Place your SQL migration files in a `migrations/` directory:

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_add_email_column.up.sql
└── 000002_add_email_column.down.sql
```

Migrations are automatically applied on startup when using `database.WithMigrations()`.

## 🤝 Contributing

This module is part of the MapBot project. See the main [MapBot repository](https://github.com/pixime/mapbot) for contribution guidelines.

## 📄 License

[MIT License](LICENSE)

## 🔗 Related Projects

- [mapbot-ai](https://github.com/pixime/mapbot-ai) - AI-powered backend
- [french-admin-etl](https://github.com/pixime/french-admin-etl) - ETL pipeline
- [mcp-sum-server](https://github.com/pixime/mcp-sum-server) - MCP server

## 📞 Support

For issues and questions, please use the [GitHub issue tracker](https://github.com/pixime/mapbot-shared/issues).
