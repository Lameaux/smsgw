package models

import (
	"math"
	"time"

	"euromoby.com/smsgw/internal/logger"
	"github.com/google/uuid"
)

func NewUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		logger.Fatal(err)
	}
	return id.String()
}

func TimeNow() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}

func CalculateNextAttemptTime(counter int) time.Time {
	return TimeNow().Add(time.Duration(30*math.Pow(2, float64(counter))) * time.Second)
}
