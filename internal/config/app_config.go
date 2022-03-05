package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	"euromoby.com/smsgw/internal/httpclient"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/utils"
)

type AppConfig struct {
	Port    string
	AppName string
	Version string

	DBPool     *pgxpool.Pool
	HTTPClient *httpclient.HTTPClient

	Merchants   map[string]string
	WorkerSleep time.Duration

	WaitTimeout       time.Duration
	ConnectionTimeout time.Duration
	TLSTimeout        time.Duration
	ReadTimeout       time.Duration
}

const (
	dbPingTimeout = 2 * time.Second
)

const defaultWorkerSleep = 5

const (
	TestAPIKey     = "test-api-key"
	TestMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29aa"
)

func defaultAppConfig(env string) *AppConfig {
	logger.Infow("loading env configuration", "env", env)

	err := godotenv.Load(".env."+env, ".env")
	if err != nil {
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
	app.configureMerchants()
	app.configurePGXPool(databaseURI)

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
	appConfig := defaultAppConfig("test")

	// Add API Key for unit tests
	appConfig.Merchants[TestAPIKey] = TestMerchantID

	return appConfig
}

func (app *AppConfig) configureMerchants() {
	app.Merchants = map[string]string{
		"postman-api-key": "d70c94da-dac4-4c0c-a6db-97f1740f29a8",
		"apikey1":         "d70c94da-dac4-4c0c-a6db-97f1740f29a9",
	}
}

func (app *AppConfig) configurePGXPool(uri string) {
	logger.Infow("connecting to db", "database_uri", uri)

	pool, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
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

func (app *AppConfig) CloseDBPool() {
	logger.Infow("closing db pool")
	app.DBPool.Close()
}

func (app *AppConfig) Shutdown() {
	app.CloseDBPool()
}
