package models

import (
	"errors"
	"testing"
)

const (
	input    = "+42012345678"
	expected = MSISDN(42012345678)
)

func TestNormalizeMSISDN(t *testing.T) {
	tests := []struct {
		msisdn     string
		normalized MSISDN
		err        error
	}{
		{"+42012345678", 42012345678, nil},
		{"0042012345678", 42012345678, nil},
		{"42012345678", 42012345678, nil},
		{"", 0, ErrInvalidMSISDN},
		{"123", 0, ErrInvalidMSISDN},
		{"abcd", 0, ErrInvalidMSISDN},
		{"00+42012345678", 0, ErrInvalidMSISDN},
		{"0+042012345678", 0, ErrInvalidMSISDN},
	}

	for _, tt := range tests {
		normalized, err := NormalizeMSISDN(tt.msisdn)
		if normalized != tt.normalized || !errors.Is(err, tt.err) {
			t.Errorf("Invalid result for %s. Expected: %d, %s. Got: %d %s.", tt.msisdn, tt.normalized, tt.err, normalized, err)
		}
	}
}

func BenchmarkNormalizeMSISDN(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			normalized, err := NormalizeMSISDN(input)
			if normalized != expected {
				b.Errorf("Invalid result for %s. Expected: %d, %s. Got: %d %s.", input, expected, err, normalized, err)
			}
		}
	})
}

func BenchmarkNormalizeMSISDNRegex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			normalized, err := NormalizeMSISDNRegex(input)
			if normalized != expected {
				b.Errorf("Invalid result for %s. Expected: %d, %s. Got: %d %s.", input, expected, err, normalized, err)
			}
		}
	})
}
