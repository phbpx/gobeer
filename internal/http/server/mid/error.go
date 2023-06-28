package mid

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	switch {
	case isFieldError(err):
		c.JSON(http.StatusBadRequest, fieldErrorResponse(err))
	case errors.Is(err, beers.ErrAlreadyExists):
		c.JSON(http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, beers.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, beers.ErrInvalidID):
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
	}
}

func isFieldError(err error) bool {
	var fe validator.ValidationErrors
	return errors.As(err, &fe)
}

func fieldErrorResponse(err error) errorResponse {
	mFieldErrors := map[string]string{}
	for _, e := range err.(validator.ValidationErrors) {
		mFieldErrors[e.Field()] = e.Tag()
	}
	return errorResponse{Error: "invalid request", Fields: mFieldErrors}
}
