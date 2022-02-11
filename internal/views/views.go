package views

import (
	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorJSON(c *gin.Context, code int, err error) {
	logger.Errorw("unhandled error", "error", err)
	c.JSON(code, ErrorResponse{Error: err.Error()})
}
