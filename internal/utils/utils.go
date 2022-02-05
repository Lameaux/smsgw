package utils

import (
	"os"
)

// GetEnv returns ENV variable or fallbacks to default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
