package dbtest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/phbpx/gobeer/internal/storage/postgres"
	"github.com/phbpx/gobeer/kit/docker"
)

// StartDB starts a database instance.
func StartDB() (*docker.Container, error) {
	image := "postgres:14-alpine"
	port := "5432"
	args := []string{
		"-e", "POSTGRES_PASSWORD=postgres",
		"-e", "POSTGRES_DB=testdb",
	}

	return docker.StartContainer(image, port, args...)
}

// StopDB stops a running database instance.
func StopDB(c *docker.Container) {
	docker.StopContainer(c.ID)
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, c *docker.Container, dbName string) (*sql.DB, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgres.Open(postgres.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       dbName,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	if err := postgres.StatusCheck(ctx, db); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	t.Log("Database ready")

	t.Log("Migrate and seed database ...")

	if err := postgres.RunMigrations(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	t.Log("Ready for testing ...")

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
	}

	return db, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *sql.DB
	Teardown func()
	t        *testing.T
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, c *docker.Container, dbName string) *Test {
	db, teardown := NewUnit(t, c, dbName)

	test := Test{
		DB:       db,
		t:        t,
		Teardown: teardown,
	}

	return &test
}
