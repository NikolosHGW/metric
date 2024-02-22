package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/db"
	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"go.uber.org/zap"
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

	err := db.InitDB(config.GetDBConnection())
	if err != nil {
		logger.Log.Info("init db", zap.Error(err))
	}
	if db.DB != nil {
		defer db.DB.Close()
	}

	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)
	handler := handlers.NewHandler(metricService, logger.Log)
	diskStrg := storage.NewDiskStorage(strg, logger.Log, config.GetFileStoragePath())
	diskService := services.NewDiskService(diskStrg, config.GetStoreInterval(), config.GetRestore())
	diskService.FillMetricStorage()
	go diskService.CollectMetrics()

	r := routes.InitRouter(handler)

	logger.Log.Info("Running server", zap.String("address", config.Address))

	return http.ListenAndServe(config.Address, r)
}
