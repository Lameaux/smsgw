package utils

import (
	"os"
	"strconv"

	"euromoby.com/smsgw/internal/logger"
)

// GetEnv returns ENV variable or fallbacks to default.
func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	logger.Fatalw("missing env variable", "key", key)

	return ""
}

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64) //nolint:gomnd
}

func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64) //nolint:gomnd
}

func FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10) //nolint:gomnd
}

func FormatUint64(u uint64) string {
	return strconv.FormatUint(u, 10) //nolint:gomnd
}
