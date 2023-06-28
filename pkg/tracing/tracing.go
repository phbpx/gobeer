// Package tracing knows how to deal with tracing.
package tracing

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// NewTracerProvider returns a new TracerProvider configured with the given options.
func NewTracerProvider(service, reporterURL string, probability float64) (*tracesdk.TracerProvider, error) {
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

	// Setup global tracer provider.
	otel.SetTracerProvider(tp)

	// Setup global text map propagator.
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	return tp, nil
}
