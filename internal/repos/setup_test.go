package repos

import (
	"os"
	"testing"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/utils"
)

var TestAppConfig *config.AppConfig

func TestMain(m *testing.M) {
	TestAppConfig = config.NewTestAppConfig()
	defer TestAppConfig.Shutdown()

	utils.CleanupDatabase(TestAppConfig.DBPool)

	os.Exit(m.Run())
}
