package testutils

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/pixime-net/mapbot-shared/config"
	"github.com/pixime-net/mapbot-shared/database"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testdbName     = "testdb"
	testdbUser     = "testuser"
	testdbPassword = "testpass"
)

// PostgresContainer embeds the testcontainers PostgresContainer and adds the ConnectionString field.
type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// SetupPostGISContainer creates a PostGIS container with testcontainers
func SetupPostGISContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgis/postgis:15-3.3",
		postgres.WithDatabase(testdbName),
		postgres.WithUsername(testdbUser),
		postgres.WithPassword(testdbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

// TestManagerOptions sets behavior of SetupTestManager
type TestManagerOptions struct {
	// MigrationsPath path to migration files (default: "/migrations" at project root)
	MigrationsPath string
	// SkipMigrations disables automatic migrations
	SkipMigrations bool
	// MigrationsSchema schema for migrations table (default: "public")
	MigrationsSchema string
	// MigrationsTable table name for migrations (default: "schema_migrations")
	MigrationsTable string
}

// getMigrationsPath returns the default migrations absolute path relative to caller
func getMigrationsPath() string {
	_, filename, _, ok := runtime.Caller(2) // Caller of SetupTestManager
	if !ok {
		panic("failed to get current file path")
	}
	// Go up to project root and look for migrations/
	dir := filepath.Dir(filename)
	for i := 0; i < 10; i++ { // Max 10 levels up
		migrationsPath := filepath.Join(dir, "migrations")
		if _, err := filepath.Abs(migrationsPath); err == nil {
			return migrationsPath
		}
		dir = filepath.Join(dir, "..")
	}
	return "./migrations"
}

// SetupTestManager creates a Manager with testcontainer and migrations
func SetupTestManager(t testing.TB, opts ...TestManagerOptions) (*database.Manager, func()) {
	t.Helper()

	// Default options
	options := TestManagerOptions{
		MigrationsPath:   getMigrationsPath(),
		SkipMigrations:   false,
		MigrationsSchema: "public",
		MigrationsTable:  "schema_migrations",
	}
	if len(opts) > 0 {
		if opts[0].MigrationsPath != "" {
			options.MigrationsPath = opts[0].MigrationsPath
		}
		options.SkipMigrations = opts[0].SkipMigrations
		if opts[0].MigrationsSchema != "" {
			options.MigrationsSchema = opts[0].MigrationsSchema
		}
		if opts[0].MigrationsTable != "" {
			options.MigrationsTable = opts[0].MigrationsTable
		}
	}

	ctx := context.Background()
	pgContainer, err := SetupPostGISContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to start PostGIS container: %v", err)
	}

	// Extract connection information from the container
	host, err := pgContainer.Host(ctx)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Create config with individual parameters
	cfg := config.NewPostgresDatabase(
		host,
		port.Int(),
		testdbName,
		testdbUser,
		testdbPassword,
	)

	// Specific configuration for tests
	cfg.MaxOpenConns = 5
	cfg.MaxIdleConns = 2
	cfg.PingTimeout = 10

	var databaseManager *database.Manager

	// Create Manager with or without migrations
	if !options.SkipMigrations {
		databaseManager, err = database.NewManager(cfg,
			database.WithMigrationsCustomSchema(
				options.MigrationsPath,
				options.MigrationsSchema,
				options.MigrationsTable,
			),
		)
		if err != nil {
			_ = pgContainer.Terminate(ctx)
			t.Fatalf("Failed to create database manager with migrations: %v", err)
		}
		t.Logf("Migrations applied from %s", options.MigrationsPath)
	} else {
		databaseManager, err = database.NewManager(cfg)
		if err != nil {
			_ = pgContainer.Terminate(ctx)
			t.Fatalf("Failed to create database manager: %v", err)
		}
	}

	cleanup := func() {
		_ = databaseManager.Close()
		if termErr := pgContainer.Terminate(ctx); termErr != nil {
			t.Logf("Failed to terminate container: %v", termErr)
		}
	}

	t.Logf("Using PostGIS testcontainer with Manager - %s", databaseManager.PublicConnectionString())
	return databaseManager, cleanup
}
