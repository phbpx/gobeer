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
	"github.com/phbpx/gobeer/internal/http/rest"
	"github.com/phbpx/gobeer/internal/storage/postgres"
	"github.com/phbpx/gobeer/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const service = "gobeer-api"

func main() {
	ctx := context.Background()
	log := logger.New(os.Stdout, logger.LevelInfo, service)

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
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
		Tracing struct {
			ReporterURI string  `conf:"default:http://localhost:14268/api/traces"`
			Probability float64 `conf:"default:0.5"`
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

	// -------------------------------------------------------------------------
	// Database Support

	// Create connectivity to the database.
	log.Info(ctx, "startup", "status", "initializing database support", "host", cfg.DB.Host)

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

	// -------------------------------------------------------------------------
	// Update the schema, if needed.

	log.Info(ctx, "startup", "status", "updating database schema", "database", cfg.DB.Name, "host", cfg.DB.Host)

	if err := postgres.RunMigrations(context.Background(), db, log); err != nil {
		log.Info(ctx, "shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		return fmt.Errorf("migrating db: %w", err)
	}

	// -------------------------------------------------------------------------
	// Start Tracing Support

	log.Info(ctx, "startup", "status", "initializing OT/Jaeger tracing support")

	traceProvider, err := startTracing(service, cfg.Tracing.ReporterURI, cfg.Tracing.Probability)
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer traceProvider.Shutdown(context.Background())

	tracer := otel.Tracer("gobeer-api")

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing http server")

	// Create handler.
	h := rest.NewHandler(rest.Config{
		Log:    log,
		Tracer: tracer,
		DB:     db,
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
		log.Info(ctx, "startup", "status", "http router started", "host", srv.Addr)
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
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

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

// ============================================================================

func startTracing(service, reporterURL string, probability float64) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(reporterURL)))
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(probability)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp,
			tracesdk.WithMaxExportBatchSize(tracesdk.DefaultMaxExportBatchSize),
			tracesdk.WithBatchTimeout(tracesdk.DefaultScheduleDelay*time.Millisecond),
			tracesdk.WithMaxExportBatchSize(tracesdk.DefaultMaxExportBatchSize),
		),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("exporter", "jaeger"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
