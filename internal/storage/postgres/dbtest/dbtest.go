package dbtest

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/phbpx/gobeer/internal/storage/postgres"
	"github.com/phbpx/gobeer/kit/docker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// New creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func New(t *testing.T, c *docker.Container) (*sql.DB, *zap.SugaredLogger, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var buf bytes.Buffer

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buf)
	log := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)).
		Sugar()

	db, err := postgres.Open(postgres.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "testdb",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	if err := postgres.StatusCheck(ctx, db, log); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	t.Log("Database ready")

	t.Log("Update database schema ...")

	if err := postgres.RunMigrations(ctx, db, log); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	t.Log("Database schema updated")

	t.Log("Ready for testing ...")

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()

		log.Sync()

		writer.Flush()
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return db, log, teardown
}
