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

func TestNormalizeMSISDN(t *testing.T) {
	tests := []struct {
		msisdn     string
		normalized string
		err        error
	}{
		{"+42012345678", "42012345678", nil},
		{"0042012345678", "42012345678", nil},
		{"42012345678", "42012345678", nil},
		{"", "", ErrInvalidMSISDN},
		{"123", "", ErrInvalidMSISDN},
		{"abcd", "", ErrInvalidMSISDN},
		{"00+42012345678", "", ErrInvalidMSISDN},
	}

	for _, tt := range tests {
		normalized, err := NormalizeMSISDN(tt.msisdn)
		if normalized != tt.normalized || err != tt.err {
			t.Errorf("Invalid result for %s. Expected: %s, %s. Got: %s %s.", tt.msisdn, tt.normalized, tt.err, normalized, err)
		}
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
