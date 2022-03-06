package handlers

import (
	"os"
	"testing"

	"euromoby.com/smsgw/internal/config"
)

var TestAppConfig *config.AppConfig //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	TestAppConfig = config.NewTestAppConfig()

	os.Exit(func() int {
		defer TestAppConfig.Shutdown()

		return m.Run()
	}())
}
