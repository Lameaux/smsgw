package models

import (
	"testing"
	"time"
)

func Test_CalculateNextAttemptTime(t *testing.T) {
	now := TimeNow()
	next := CalculateNextAttemptTime(0)

	if diff := next.Sub(now); diff < 30*time.Second {
		t.Errorf("Invalid interval. Expected: 30. Got %v", diff)
	}
}
