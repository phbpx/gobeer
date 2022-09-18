package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/http/rest/mid"
	"github.com/phbpx/gobeer/internal/listing"
	"github.com/phbpx/gobeer/internal/reviewing"
)

// Config holds the dependencies for the handler.
type Config struct {
	Adding    *adding.Service
	Reviewing *reviewing.Service
	Listing   *listing.Service
}

// Handler is the HTTP handler for the REST API.
type Handler struct {
	adding    *adding.Service
	reviewing *reviewing.Service
	listing   *listing.Service
}

// NewHandler creates a new Handler.
func NewHandler(cfg Config) *Handler {
	return &Handler{
		adding:    cfg.Adding,
		reviewing: cfg.Reviewing,
		listing:   cfg.Listing,
	}
}

// Router returns the gin router.
func (h *Handler) Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Add middlewares.
	r.Use(gin.Recovery())
	r.Use(mid.ErrorHandler())

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
func (h *Handler) addBeer(c *gin.Context) {
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
func (h *Handler) listBeers(c *gin.Context) {
	ctx := c.Request.Context()

	bs, err := h.listing.ListBeers(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, bs)
}

// addReview is the HTTP handler for the POST /beers/:id/reviews endpoint.
func (h *Handler) addReview(c *gin.Context) {
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
func (h *Handler) listReviews(c *gin.Context) {
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
