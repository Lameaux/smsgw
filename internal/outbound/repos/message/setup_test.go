package message

import (
	"os"
	"testing"

	"euromoby.com/smsgw/internal/config"
)

var TestApp *config.App //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	TestApp = config.NewTestApp()

	os.Exit(func() int {
		defer TestApp.Config.Shutdown()

		return m.Run()
	}())
}
