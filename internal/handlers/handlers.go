package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func queryParamIntDefault(c *gin.Context, param string, def int) (int, error) {
	value := c.Query(param)
	if value == "" {
		return def, nil
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return def, err
	}

	return i, nil
}

func queryParamTime(c *gin.Context, param string) (*time.Time, error) {
	value := c.Query(param)
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
