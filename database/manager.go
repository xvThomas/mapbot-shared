package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pixime/mapbot-shared/config"

	"github.com/jackc/pgx/v5/pgxpool"
	// Import the pgx driver for database/sql
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DatabaseManager manages connections to the database and provides utility methods for health checks and stats
type DatabaseManager struct {
	db     *sql.DB
	config *config.PostgresDatabase
	pool   *pgxpool.Pool
}

// DatabaseManagerOption is a configuration function
type DatabaseManagerOption func(*DatabaseManager) error

// WithMigrations is an option to automatically run migrations
func WithMigrations(migrationsPath string) DatabaseManagerOption {
	return func(dm *DatabaseManager) error {
		return RunMigrations(dm.db, migrationsPath, "public", "schema_migrations")
	}
}

// WithMigrationsCustomSchema is an option to run migrations with custom schema and table
func WithMigrationsCustomSchema(migrationsPath, schemaName, tableName string) DatabaseManagerOption {
	return func(dm *DatabaseManager) error {
		return RunMigrations(dm.db, migrationsPath, schemaName, tableName)
	}
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(cfg *config.PostgresDatabase, opts ...DatabaseManagerOption) (*DatabaseManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config cannot be nil")
	}

	// Validation of required fields
	if cfg.Host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if cfg.Database == "" {
		return nil, fmt.Errorf("database name cannot be empty")
	}
	if cfg.User == "" {
		return nil, fmt.Errorf("user cannot be empty")
	}

	connectionString := cfg.ConnectionString()
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure the connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.PingTimeout)*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to database %s within %ds: %w",
			cfg.PublicConnectionString(), cfg.PingTimeout, err)
	}

	// Create pgxpool for advanced features (batch operations, better performance)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error parsing pool config: %w", err)
	}

	// Configure pool with similar settings
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	dm := &DatabaseManager{
		db:     db,
		config: cfg,
		pool:   pool,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(dm); err != nil {
			_ = dm.Close()
			return nil, fmt.Errorf("failed to apply database option: %w", err)
		}
	}

	return dm, nil
}

// GetDB returns the *sql.DB instance
func (dm *DatabaseManager) GetDB() *sql.DB {
	return dm.db
}

// GetPool returns the *pgxpool.Pool instance for advanced operations
func (dm *DatabaseManager) GetPool() *pgxpool.Pool {
	return dm.pool
}

// GetConfig returns the database configuration
func (dm *DatabaseManager) GetConfig() *config.PostgresDatabase {
	return dm.config
}

// Ping tests the connection to the database
func (dm *DatabaseManager) Ping(ctx context.Context) error {
	return dm.db.PingContext(ctx)
}

// Stats returns the connection pool statistics
func (dm *DatabaseManager) Stats() sql.DBStats {
	return dm.db.Stats()
}

// Close closes all connections
func (dm *DatabaseManager) Close() error {
	var dbErr, poolErr error
	if dm.pool != nil {
		dm.pool.Close()
	}
	if dm.db != nil {
		dbErr = dm.db.Close()
	}
	if dbErr != nil {
		return dbErr
	}
	return poolErr
}

// Health checks the health status of the database
func (dm *DatabaseManager) Health(ctx context.Context) error {
	if err := dm.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	stats := dm.Stats()
	if stats.OpenConnections >= dm.config.MaxOpenConns {
		return fmt.Errorf("database connection pool exhausted: %d/%d connections",
			stats.OpenConnections, dm.config.MaxOpenConns)
	}

	return nil
}

// PublicConnectionString returns the configuration information as a string without the password
func (dm *DatabaseManager) PublicConnectionString() string {
	return dm.config.PublicConnectionString()
}
