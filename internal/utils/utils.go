package utils

import (
	"errors"
	"math"
	"os"
	"regexp"
	"time"

	"euromoby.com/smsgw/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	msisdnRegex      = regexp.MustCompile(`^(\+|00)?([1-9]\d{7,14})$`)
	ErrInvalidMSISDN = errors.New("invalid msisdn format")
)

func ErrorJSON(c *gin.Context, code int, err error) {
	c.JSON(code, ErrorResponse{Error: err.Error()})
}

// GetEnv returns ENV variable or fallbacks to default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// NewUUID returns a new UUID as string
func NewUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		logger.Fatal(err)
	}
	return id.String()
}

func Now() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}

func NormalizeMSISDN(msisdn string) (string, error) {
	match := msisdnRegex.FindStringSubmatch(msisdn)
	if match == nil {
		return "", ErrInvalidMSISDN
	}

	return match[2], nil
}

func CalculateNextAttemptTime(counter int) time.Time {
	return Now().Add(time.Duration(30*math.Pow(2, float64(counter))) * time.Second)
}
