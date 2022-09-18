package mid

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phbpx/gobeer/internal/beers"
)

// errorReponse is the JSON response for an error.
type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// ErrorHandler is the middleware for handling errors.
func ErrorHandler() gin.HandlerFunc {
	return errorHandler
}

func errorHandler(c *gin.Context) {
	c.Next()

	// If no errors, just return.
	if len(c.Errors) == 0 {
		return
	}

	// Get the last error.
	err := c.Errors.Last().Err

	switch err {
	case beers.ErrAlreadyExists:
		c.JSON(http.StatusConflict, errorResponse{Error: err.Error()})
	case beers.ErrNotFound:
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
	case beers.ErrInvalidID:
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
	}
}
