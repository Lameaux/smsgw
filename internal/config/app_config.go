package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	"euromoby.com/smsgw/internal/auth"
	"euromoby.com/smsgw/internal/billing"
	"euromoby.com/smsgw/internal/httpclient"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/utils"
)

type AppConfig struct {
	Env string

	Port    string
	AppName string
	Version string

	DBPool     *pgxpool.Pool
	HTTPClient *httpclient.HTTPClient

	WorkerSleep time.Duration

	WaitTimeout       time.Duration
	ConnectionTimeout time.Duration
	TLSTimeout        time.Duration
	ReadTimeout       time.Duration

	Auth    auth.Auth
	Billing billing.Billing
}

const (
	dbPingTimeout = 2 * time.Second

	defaultWorkerSleep = 5
)

func defaultAppConfig(env string) *AppConfig {
	logger.Infow("loading env configuration", "env", env)

	if err := godotenv.Load(".env."+env, ".env"); err != nil {
		logger.Fatalw("failed to load env", "env", env, "error", err)

		return nil
	}

	port := utils.GetEnv("PORT")
	databaseURI := utils.GetEnv("DATABASE_URI")

	waitTimeout, err := strconv.Atoi(utils.GetEnv("WAIT_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	connectionTimeout, err := strconv.Atoi(utils.GetEnv("CONNECTION_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	tlsTimeout, err := strconv.Atoi(utils.GetEnv("TLS_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	readTimeout, err := strconv.Atoi(utils.GetEnv("READ_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	app := &AppConfig{
		Env: env,

		AppName: "smsgw",
		Version: "0.1",
		Port:    port,

		WaitTimeout:       time.Duration(waitTimeout) * time.Second,
		ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
		TLSTimeout:        time.Duration(tlsTimeout) * time.Second,
		ReadTimeout:       time.Duration(readTimeout) * time.Second,

		WorkerSleep: time.Duration(defaultWorkerSleep) * time.Second,
	}

	app.configureHTTPClient()
	app.configurePGXPool(databaseURI)

	app.configureAuth()
	app.configureBilling()

	return app
}

func NewAppConfig() *AppConfig {
	env := os.Getenv("SMSGW_ENV")
	if env == "" {
		env = "development"
	}

	return defaultAppConfig(env)
}

func NewTestAppConfig() *AppConfig {
	utils.SetWorkingDir()

	app := defaultAppConfig("test")

	utils.CleanupDatabase(app.DBPool)

	return app
}

func (app *AppConfig) configurePGXPool(uri string) {
	logger.Infow("connecting to db", "database_uri", uri)

	pool, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal(err)
	}

	app.DBPool = pool
}

func (app *AppConfig) configureHTTPClient() {
	app.HTTPClient = httpclient.NewBuilder().
		ConnectionTimeout(app.ConnectionTimeout).
		TLSTimeout(app.TLSTimeout).
		ReadTimeout(app.ReadTimeout).
		Build()
}

func (app *AppConfig) configureAuth() {
	if app.Env == "test" {
		app.Auth = auth.NewTestAuth()
	} else {
		app.Auth = auth.NewStubAuth()
	}
}

func (app *AppConfig) configureBilling() {
	if app.Env == "test" {
		app.Billing = billing.NewTestBilling()
	} else {
		app.Billing = billing.NewStubBilling()
	}
}

func (app *AppConfig) CloseDBPool() {
	logger.Infow("closing db pool")
	app.DBPool.Close()
}

func (app *AppConfig) Shutdown() {
	app.CloseDBPool()
}
