// Package handler ...
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var propagator = otel.GetTextMapPropagator()

// New creates a new http.Handler.
func New(tracer trace.Tracer) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/users/{userID}/notify", func(w http.ResponseWriter, r *http.Request) {
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(ctx, "email-api")
		defer span.End()

		vars := mux.Vars(r)
		userID := vars["userID"]
		span.SetAttributes(attribute.String("user_id", userID))

		doSomething(ctx, tracer)

		time.Sleep(3 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	return router
}

func doSomething(ctx context.Context, tracer trace.Tracer) {
	_, span := tracer.Start(ctx, "external-service")
	defer span.End()

	time.Sleep(10 * time.Millisecond)
}
