package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/db"
	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

const defaultTagValue = "N/A"

var (
	buildVersion = defaultTagValue
	buildDate    = defaultTagValue
	buildCommit  = defaultTagValue
)

func main() {
	if err := run(); err != nil {
		log.Fatal(fmt.Errorf("server/main run_ListenAndServe: %w", err))
	}
}

func run() error {
	config := NewConfig()

	if err := logger.Initialize(config.LogLevel); err != nil {
		return err
	}

	database, err := db.InitDB(config.GetDBConnection())
	if err != nil {
		logger.Log.Info("init db", zap.Error(err))
	}
	if database != nil {
		defer func() {
			err := database.Close()
			if err != nil {
				logger.Log.Info("err close database", zap.Error(err))
			}
		}()
	}

	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)
	if database != nil {
		databaseStrg := storage.NewDBStorage(database, logger.Log)
		metricService = services.NewMetricService(databaseStrg)
	}
	handler := handlers.NewHandler(metricService, logger.Log)
	diskStrg := storage.NewDiskStorage(strg, logger.Log, config.GetFileStoragePath())
	diskService := services.NewDiskService(diskStrg, config.GetStoreInterval(), config.GetRestore())
	diskService.FillMetricStorage()
	go diskService.CollectMetrics()

	hashMiddleware := middlewares.NewHashMiddleware(config.GetKey())

	r := routes.InitRouter(handler, hashMiddleware)

	fmt.Println(
		"Build version: ", buildVersion, "\n",
		"Build date: ", buildDate, "\n",
		"Build commit: ", buildCommit,
	)

	logger.Log.Info("Running server", zap.String("address", config.Address))

	return http.ListenAndServe(config.Address, r)
}
