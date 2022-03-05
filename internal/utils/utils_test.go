package utils

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	if GetEnv("GOROOT") == "" {
		t.Fatalf("GOROOT is missing")
	}
}
