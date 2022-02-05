package views

import (
	"euromoby.com/smsgw/internal/logger"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorJSON(c *gin.Context, code int, err error) {
	logger.Errorw("unhandled error", "error", err)
	c.JSON(code, ErrorResponse{Error: err.Error()})
}
