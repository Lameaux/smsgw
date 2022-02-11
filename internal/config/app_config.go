package config

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

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
	defaultPort = "8080"

	defaultWaitTimeout       = "15"
	defaultConnectionTimeout = "5"
	defaultTLSTimeout        = "5"
	defaultReadTimeout       = "10"

	dbPingTimeout = 2 * time.Second
)

const defaultWorkerSleep = 5

const (
	TestAPIKey     = "test-api-key"
	TestMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29aa"
)

func defaultAppConfig(databaseURI string) *AppConfig {
	port := utils.GetEnv("PORT", defaultPort)

	waitTimeout, err := strconv.Atoi(utils.GetEnv("WAIT_TIMEOUT", defaultWaitTimeout))
	if err != nil {
		logger.Fatal(err)
	}

	connectionTimeout, err := strconv.Atoi(utils.GetEnv("CONNECTION_TIMEOUT", defaultConnectionTimeout))
	if err != nil {
		logger.Fatal(err)
	}

	tlsTimeout, err := strconv.Atoi(utils.GetEnv("TLS_TIMEOUT", defaultTLSTimeout))
	if err != nil {
		logger.Fatal(err)
	}

	readTimeout, err := strconv.Atoi(utils.GetEnv("READ_TIMEOUT", defaultReadTimeout))
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
	databaseURI := utils.GetEnv("DATABASE_URI", "postgres://root:heslo@localhost:5432/smsgw_test?&pool_max_conns=10")
	logger.Infow("Connecting to database", "database_uri", databaseURI)

	appConfig := defaultAppConfig(databaseURI)
	logger.Infow("Starting", "app", appConfig.AppName, "version", appConfig.Version)

	return appConfig
}

func NewTestAppConfig() *AppConfig {
	testDatabaseURI := utils.GetEnv("DATABASE_URI", "postgres://root:heslo@localhost:5432/smsgw_test?&pool_max_conns=10")
	appConfig := defaultAppConfig(testDatabaseURI)

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
	app.DBPool.Close()
}

func (app *AppConfig) Shutdown() {
	logger.Infow("closing DB pool")
	app.CloseDBPool()
}
