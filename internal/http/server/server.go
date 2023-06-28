// Package server contains the http server.
package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/http/server/mid"
	"github.com/phbpx/gobeer/internal/listing"
	"github.com/phbpx/gobeer/internal/reviewing"
	"github.com/phbpx/gobeer/internal/storage/postgres"
	"github.com/phbpx/gobeer/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

// Config holds the dependencies for the handler.
type Config struct {
	Log    *logger.Logger
	Tracer trace.Tracer
	DB     *sql.DB
}

// Server is the HTTP Server for the REST API.
type Server struct {
	log       *logger.Logger
	tracer    trace.Tracer
	adding    *adding.Service
	reviewing *reviewing.Service
	listing   *listing.Service
}

// New creates a new Server.
func New(cfg Config) *Server {
	storage := postgres.NewStorage(cfg.DB)
	addingSrv := adding.NewService(storage)
	reviewingSrv := reviewing.NewService(storage)
	listingSrv := listing.NewService(storage)

	return &Server{
		log:       cfg.Log,
		tracer:    cfg.Tracer,
		adding:    addingSrv,
		reviewing: reviewingSrv,
		listing:   listingSrv,
	}
}

// Router returns the gin router.
func (h *Server) Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Add middlewares.
	r.Use(
		gin.Recovery(),
		mid.Tracing(h.tracer),
		mid.Logger(h.log),
		mid.ErrorHandler(),
	)

	// app routes.
	r.POST("/beers", h.addBeer)
	r.GET("/beers", h.listBeers)
	r.POST("/beers/:id/reviews", h.addReview)
	r.GET("/beers/:id/reviews", h.listReviews)

	// debug routes.
	r.GET("/debug/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}

// addBeer is the HTTP handler for the POST /beers endpoint.
func (h *Server) addBeer(c *gin.Context) {
	ctx := c.Request.Context()

	var nb adding.NewBeer
	if err := c.ShouldBindJSON(&nb); err != nil {
		c.Error(err)
		return
	}

	b, err := h.adding.AddBeer(ctx, nb)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, b)
}

// listBeers is the HTTP handler for the GET /beers endpoint.
func (h *Server) listBeers(c *gin.Context) {
	ctx := c.Request.Context()

	bs, err := h.listing.ListBeers(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	if len(bs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, bs)
}

// addReview is the HTTP handler for the POST /beers/:id/reviews endpoint.
func (h *Server) addReview(c *gin.Context) {
	ctx := c.Request.Context()

	var nr reviewing.NewReview
	if err := c.ShouldBindJSON(&nr); err != nil {
		c.Error(err)
		return
	}

	beerID := c.Param("id")
	nr.BeerID = beerID

	bs, err := h.reviewing.CreateReview(ctx, nr)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, bs)
}

// listReviews is the HTTP handler for the GET /beers/:id/reviews endpoint.
func (h *Server) listReviews(c *gin.Context) {
	ctx := c.Request.Context()
	beerID := c.Param("id")

	r, err := h.listing.ListReviews(ctx, beerID)
	if err != nil {
		c.Error(err)
		return
	}

	if len(r) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, r)
}
