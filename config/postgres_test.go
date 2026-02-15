package config

import (
	"strings"
	"testing"
)

// TestNewPostgresDatabase tests creating a new PostgresDatabase with default values
func TestNewPostgresDatabase(t *testing.T) {
	db := NewPostgresDatabase("localhost", 5432, "testdb", "testuser", "testpass")

	if db.Host != "localhost" {
		t.Errorf("Host = %s, want localhost", db.Host)
	}
	if db.Port != 5432 {
		t.Errorf("Port = %d, want 5432", db.Port)
	}
	if db.Database != "testdb" {
		t.Errorf("Database = %s, want testdb", db.Database)
	}
	if db.User != "testuser" {
		t.Errorf("User = %s, want testuser", db.User)
	}
	if db.Password != "testpass" {
		t.Errorf("Password = %s, want testpass", db.Password)
	}

	// Check default values
	if db.SSLMode != "disable" {
		t.Errorf("SSLMode = %s, want disable", db.SSLMode)
	}
	if db.MaxOpenConns != 25 {
		t.Errorf("MaxOpenConns = %d, want 25", db.MaxOpenConns)
	}
	if db.MaxIdleConns != 5 {
		t.Errorf("MaxIdleConns = %d, want 5", db.MaxIdleConns)
	}
	if db.ConnMaxLifetime != 5 {
		t.Errorf("ConnMaxLifetime = %d, want 5", db.ConnMaxLifetime)
	}
	if db.ConnMaxIdleTime != 30 {
		t.Errorf("ConnMaxIdleTime = %d, want 30", db.ConnMaxIdleTime)
	}
	if db.PingTimeout != 10 {
		t.Errorf("PingTimeout = %d, want 10", db.PingTimeout)
	}
}

// TestConnectionString tests connection string generation
func TestConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		db       *PostgresDatabase
		expected string
	}{
		{
			name: "basic connection",
			db: &PostgresDatabase{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				User:     "user",
				Password: "pass",
				SSLMode:  "disable",
			},
			expected: "postgres://user:pass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "with SSL enabled",
			db: &PostgresDatabase{
				Host:     "db.example.com",
				Port:     5433,
				Database: "proddb",
				User:     "admin",
				Password: "secret123",
				SSLMode:  "require",
			},
			expected: "postgres://admin:secret123@db.example.com:5433/proddb?sslmode=require",
		},
		{
			name: "with special characters in password",
			db: &PostgresDatabase{
				Host:     "localhost",
				Port:     5432,
				Database: "mydb",
				User:     "user",
				Password: "p@ss!word#123",
				SSLMode:  "disable",
			},
			expected: "postgres://user:p@ss!word#123@localhost:5432/mydb?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.db.ConnectionString()
			if got != tt.expected {
				t.Errorf("ConnectionString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestPublicConnectionString tests public connection string generation (password masked)
func TestPublicConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		db       *PostgresDatabase
		expected string
	}{
		{
			name: "password masked",
			db: &PostgresDatabase{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				User:     "user",
				Password: "secretpassword",
				SSLMode:  "disable",
			},
			expected: "postgres://user:****@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "production config",
			db: &PostgresDatabase{
				Host:     "db.prod.com",
				Port:     5433,
				Database: "proddb",
				User:     "admin",
				Password: "verysecret",
				SSLMode:  "require",
			},
			expected: "postgres://admin:****@db.prod.com:5433/proddb?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.db.PublicConnectionString()
			if got != tt.expected {
				t.Errorf("PublicConnectionString() = %q, want %q", got, tt.expected)
			}

			// Verify password is not present
			if strings.Contains(got, tt.db.Password) {
				t.Errorf("PublicConnectionString() contains password %q", tt.db.Password)
			}

			// Verify it contains the mask
			if !strings.Contains(got, "****") {
				t.Error("PublicConnectionString() should contain password mask (****)")
			}
		})
	}
}

// TestConnectionStringFormat tests the format of connection strings
func TestConnectionStringFormat(t *testing.T) {
	db := NewPostgresDatabase("host", 1234, "db", "user", "pass")

	connStr := db.ConnectionString()

	// Verify format components
	if !strings.HasPrefix(connStr, "postgres://") {
		t.Error("ConnectionString() should start with postgres://")
	}
	if !strings.Contains(connStr, "@") {
		t.Error("ConnectionString() should contain @ separator")
	}
	if !strings.Contains(connStr, "?sslmode=") {
		t.Error("ConnectionString() should contain sslmode parameter")
	}
}

// TestPublicConnectionStringDoesNotLeakPassword tests password is never leaked
func TestPublicConnectionStringDoesNotLeakPassword(t *testing.T) {
	passwords := []string{
		"simple",
		"c0mpl3x!P@ss",
		"with spaces",
		"with-dashes",
		"with_underscores",
		"",
	}

	for _, password := range passwords {
		db := &PostgresDatabase{
			Host:     "localhost",
			Port:     5432,
			Database: "db",
			User:     "user",
			Password: password,
			SSLMode:  "disable",
		}

		publicStr := db.PublicConnectionString()

		if password != "" && strings.Contains(publicStr, password) {
			t.Errorf("PublicConnectionString() leaked password %q in output: %s", password, publicStr)
		}
	}
}

// TestDefaultValues tests that default values are sensible
func TestDefaultValues(t *testing.T) {
	db := NewPostgresDatabase("host", 5432, "db", "user", "pass")

	// Connection pool defaults
	if db.MaxOpenConns <= 0 {
		t.Error("MaxOpenConns should be positive")
	}
	if db.MaxIdleConns <= 0 {
		t.Error("MaxIdleConns should be positive")
	}
	if db.MaxIdleConns > db.MaxOpenConns {
		t.Error("MaxIdleConns should not exceed MaxOpenConns")
	}

	// Timeout defaults
	if db.ConnMaxLifetime <= 0 {
		t.Error("ConnMaxLifetime should be positive")
	}
	if db.ConnMaxIdleTime <= 0 {
		t.Error("ConnMaxIdleTime should be positive")
	}
	if db.PingTimeout <= 0 {
		t.Error("PingTimeout should be positive")
	}
}

// TestModifyConfiguration tests that configuration can be modified after creation
func TestModifyConfiguration(t *testing.T) {
	db := NewPostgresDatabase("localhost", 5432, "testdb", "user", "pass")

	// Modify values
	db.MaxOpenConns = 100
	db.MaxIdleConns = 20
	db.SSLMode = "require"

	// Verify modifications
	if db.MaxOpenConns != 100 {
		t.Errorf("MaxOpenConns = %d, want 100", db.MaxOpenConns)
	}
	if db.MaxIdleConns != 20 {
		t.Errorf("MaxIdleConns = %d, want 20", db.MaxIdleConns)
	}
	if db.SSLMode != "require" {
		t.Errorf("SSLMode = %s, want require", db.SSLMode)
	}

	// Verify connection string reflects changes
	connStr := db.ConnectionString()
	if !strings.Contains(connStr, "sslmode=require") {
		t.Error("ConnectionString() should reflect SSLMode change")
	}
}
