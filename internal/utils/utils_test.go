package utils

import (
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	if GetEnv("SOME_RANDOM_DUMMY_VALUE", "fallback") != "fallback" {
		t.Fatalf("Invalid fallback value")
	}

	if GetEnv("GOROOT", "fallback") == "fallback" {
		t.Fatalf("GOROOT is missing")
	}
}

func Test_CalculateNextAttemptTime(t *testing.T) {
	now := Now()
	next := CalculateNextAttemptTime(0)
	diff := next.Sub(now)

	if diff < 30*time.Second {
		t.Errorf("Invalid interval. Expected: 30. Got %v", diff)
	}
}
