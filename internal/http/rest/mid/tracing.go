package mid

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerKey = "gin-tracer"

// Logger is a middleware that logs the request as it goes in and the response as it goes out.
func Tracing(tracer trace.Tracer) gin.HandlerFunc {
	textMapPropagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		c.Set(tracerKey, tracer)

		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		// create ctx with propagators.
		ctx := textMapPropagator.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))

		path := c.FullPath()
		method := c.Request.Method
		if path == "" {
			path = fmt.Sprintf("HTTP %s route not found", method)
		}
		spanName := fmt.Sprintf("%s %s", method, path)

		// Create a span
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPScheme(c.Request.URL.Scheme),
				semconv.HTTPMethod(method),
				semconv.HTTPURL(c.Request.URL.String()),
			),
		)
		defer span.End()

		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		span.SetAttributes(semconv.HTTPStatusCode(status))

		if status >= 400 {
			span.SetStatus(codes.Error, "")
		}

		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}
