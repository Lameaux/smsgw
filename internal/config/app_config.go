package config

import (
	"context"
	"strconv"
	"time"

	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AppConfig struct {
	Port        string
	AppName     string
	Version     string
	WaitTimeout time.Duration
	DBPool      *pgxpool.Pool
	Merchants   map[string]string
	WorkerSleep time.Duration
}

const appName = "smsgw"
const version = "0.1"

const defaultWaitTimeout = "15"
const defaultPort = "8080"

const defaultWorkerSleep = 5

const TestApiKey = "test-api-key"
const TestMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29aa"

func DefaultAppConfig() *AppConfig {
	port := utils.GetEnv("PORT", defaultPort)

	waitTimeout, err := strconv.Atoi(utils.GetEnv("WAIT_TIMEOUT", defaultWaitTimeout))
	if err != nil {
		logger.Fatal(err)
	}

	return &AppConfig{
		AppName:     appName,
		Version:     version,
		Port:        port,
		WaitTimeout: time.Duration(waitTimeout) * time.Second,
		WorkerSleep: time.Duration(defaultWorkerSleep) * time.Second,
		Merchants:   loadMerchants(),
	}

}

func NewAppConfig() *AppConfig {
	logger.Infow("Starting", "app", appName, "version", version)

	databaseURI := utils.GetEnv("DATABASE_URI", "postgres://root:heslo@localhost:5432/smsgw_test?&pool_max_conns=10")
	logger.Infow("Connecting to database", "database_uri", databaseURI)

	appConfig := DefaultAppConfig()
	appConfig.DBPool = newPGXPool(databaseURI)

	return appConfig
}

func NewTestAppConfig() *AppConfig {
	testDatabaseURI := utils.GetEnv("DATABASE_URI", "postgres://root:heslo@localhost:5432/smsgw_test?&pool_max_conns=10")
	appConfig := DefaultAppConfig()
	appConfig.DBPool = newPGXPool(testDatabaseURI)

	// Add API Key for unit tests
	appConfig.Merchants[TestApiKey] = TestMerchantID

	return appConfig
}

func loadMerchants() map[string]string {
	return map[string]string{
		"postman-api-key": "d70c94da-dac4-4c0c-a6db-97f1740f29a8",
		"apikey1":         "d70c94da-dac4-4c0c-a6db-97f1740f29a9",
	}
}

func newPGXPool(uri string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	return pool
}

func (config *AppConfig) CloseDBPool() {
	config.DBPool.Close()
}

func (config *AppConfig) Shutdown() {
	defer config.CloseDBPool()
}
