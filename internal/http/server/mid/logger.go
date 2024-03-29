package mid

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phbpx/gobeer/pkg/logger"
)

// Logger is a middleware that logs the request as it goes in and the response as it goes out.
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if len(c.Errors) > 0 {
			log.Error(ctx, "request",
				"path", fmt.Sprintf("%s %s%s", method, path, query),
				"status", status,
				"latency", latency,
				"ERROR", c.Errors,
			)
			return
		}

		log.Info(ctx, "request",
			"path", fmt.Sprintf("%s %s%s", method, path, query),
			"status", status,
			"latency", latency,
		)
	}
}
