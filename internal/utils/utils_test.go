package utils

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	if GetEnv("SOME_RANDOM_DUMMY_VALUE", "fallback") != "fallback" {
		t.Fatalf("Invalid fallback value")
	}

	if GetEnv("GOROOT", "fallback") == "fallback" {
		t.Fatalf("GOROOT is missing")
	}
}
