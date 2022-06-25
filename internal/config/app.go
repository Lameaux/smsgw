package config

import (
	"os"
	"time"

	coreconfig "euromoby.com/core/config"
	coretesting "euromoby.com/core/testing"
)

var tables = []string{ //nolint:gochecknoglobals
	"message_groups",
	"outbound_messages",
	"inbound_messages",
	"delivery_notifications",
	"keys",
	"shortcodes",
	"callbacks",
}

const (
	AppName    = "smsgw"
	AppVersion = "0.2"

	defaultWorkerSleep = 5
)

type App struct {
	Config coreconfig.AppConfig

	WorkerSleep time.Duration
}

func NewApp() *App {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = coreconfig.EnvDevelopment
	}

	app := defaultApp(env)
	return app
}

func NewTestApp() *App {
	coretesting.SetWorkingDir()

	app := defaultApp(coreconfig.EnvTest)

	coretesting.CleanupDatabase(app.Config.DBPool, tables)

	return app
}

func defaultApp(env string) *App {
	return &App{
		Config:      *coreconfig.NewAppConfig(env),
		WorkerSleep: time.Duration(defaultWorkerSleep) * time.Second,
	}
}
