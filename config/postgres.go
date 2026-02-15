package config

import "fmt"

// PostgresDatabase configuration structure for PostgreSQL connections
type PostgresDatabase struct {
	Host            string
	Port            int
	Database        string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // in minutes
	ConnMaxIdleTime int // in seconds
	PingTimeout     int // in seconds
}

// NewPostgresDatabase creates a PostgresDatabase config with sensible defaults
func NewPostgresDatabase(host string, port int, database, user, password string) *PostgresDatabase {
	return &PostgresDatabase{
		Host:            host,
		Port:            port,
		Database:        database,
		User:            user,
		Password:        password,
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5,
		ConnMaxIdleTime: 30,
		PingTimeout:     10,
	}
}

// ConnectionString returns the PostgreSQL connection string
func (p *PostgresDatabase) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Database, p.SSLMode)
}

// PublicConnectionString returns the connection string without password for logging
func (p *PostgresDatabase) PublicConnectionString() string {
	return fmt.Sprintf("postgres://%s:****@%s:%d/%s?sslmode=%s",
		p.User, p.Host, p.Port, p.Database, p.SSLMode)
}
