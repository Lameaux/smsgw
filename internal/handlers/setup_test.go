package handlers

import (
	"os"
	"testing"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/utils"
)

var TestAppConfig *config.AppConfig //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	TestAppConfig = config.NewTestAppConfig()

	utils.CleanupDatabase(TestAppConfig.DBPool)

	os.Exit(func() int {
		defer TestAppConfig.Shutdown()

		return m.Run()
	}())
}
