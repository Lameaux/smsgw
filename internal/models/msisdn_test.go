package models

import "testing"

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
	}

	for _, tt := range tests {
		normalized, err := NormalizeMSISDN(tt.msisdn)
		if normalized != tt.normalized || err != tt.err {
			t.Errorf("Invalid result for %s. Expected: %s, %s. Got: %s %s.", tt.msisdn, tt.normalized, tt.err, normalized, err)
		}
	}
}
