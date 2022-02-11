package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			cancel()

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				c.AbortWithStatus(http.StatusGatewayTimeout)
			}
		}()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
