package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/http/rest"
	"github.com/phbpx/gobeer/internal/listing"
	"github.com/phbpx/gobeer/internal/reviewing"
	"github.com/phbpx/gobeer/internal/storage/postgres"
	"github.com/phbpx/gobeer/kit/logger"
	"go.uber.org/zap"
)

func main() {
	// Construct the application logger.
	log, err := logger.New("gobeer-api")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Server struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:localhost"`
			Name         string `conf:"default:testdb"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{}

	const prefix = "GOBEER"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)

	}

	// =========================================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := postgres.Open(postgres.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		db.Close()
	}()

	// =========================================================================
	// Update the schema, if needed.

	log.Infow("startup", "status", "updating database schema", "database", cfg.DB.Name, "host", cfg.DB.Host)

	if err := postgres.RunMigrations(context.Background(), db); err != nil {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		return fmt.Errorf("migrating db: %w", err)
	}

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing http server")

	// Create dependencies for the API.
	repository := postgres.NewStorage(db)
	addingSvc := adding.NewService(repository)
	reviewingSvc := reviewing.NewService(repository)
	listingSvc := listing.NewService(repository)

	// Create handler.
	h := rest.NewHandler(rest.Config{
		Adding:    addingSvc,
		Reviewing: reviewingSvc,
		Listing:   listingSvc,
	})

	// Create a new HTTP server.
	srv := http.Server{
		Addr:         cfg.Server.APIHost,
		Handler:      h.Router(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "http router started", "host", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
