// Package postgres provides a storage implementation for PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
	"github.com/phbpx/gobeer/pkg/logger"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

//go:embed migrations
var migrations embed.FS

// Config is the required properties to use the database.
type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sql.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := otelsql.Open("postgres", u.String(), otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
		semconv.DBName(cfg.Name),
	))
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB, log *logger.Logger) error {

	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		log.Warn(ctx, "postgres: ping attempt %d failed: %v", attempts, pingError)
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// RunMigrations runs the database migrations.
func RunMigrations(ctx context.Context, db *sql.DB, log *logger.Logger) error {
	// Check if the database is ready.
	if err := StatusCheck(ctx, db, log); err != nil {
		return fmt.Errorf("db status check: %w", err)
	}
	// Load the migrations from the embedded filesystem.
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("invalid source instance: %w", err)
	}

	// Create the database driver for the migrations.
	target, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("invalid target postgres instance, %w", err)
	}

	// Create the migration instance.
	m, err := migrate.NewWithInstance("httpfs", source, "postgres", target)
	if err != nil {
		return err
	}

	// Run the migrations.
	if err := m.Up(); err != nil {
		// If the error is not a "no change" error, return it.
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
