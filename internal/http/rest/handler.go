package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/http/rest/mid"
)

// Config holds the dependencies for the handler.
type Config struct {
	Adding *adding.Service
}

// Handler is the HTTP handler for the REST API.
type Handler struct {
	adding *adding.Service
}

// NewHandler creates a new Handler.
func NewHandler(cfg Config) *Handler {
	return &Handler{
		adding: cfg.Adding,
	}
}

// Router returns the gin router.
func (h *Handler) Router() *gin.Engine {
	r := gin.New()

	// Add middlewares.
	r.Use(gin.Recovery())
	r.Use(mid.ErrorHandler())

	// Set routes.
	r.POST("/beers", h.addBeer)

	return r
}

// addBeer is the HTTP handler for the POST /beers endpoint.
func (h *Handler) addBeer(c *gin.Context) {
	ctx := c.Request.Context()

	var nb adding.NewBeer
	if err := c.BindJSON(&nb); err != nil {
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
