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
	"github.com/NikolosHGW/metric/internal/server/storage/disk"
	"github.com/NikolosHGW/metric/internal/server/storage/memory"
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

	db.InitDB()
	if db.DB != nil {
		defer db.DB.Close()
	}

	strg := memory.NewMemStorage()
	metricService := services.NewMetricService(strg)
	handler := handlers.NewHandler(metricService, logger.Log)
	diskStrg := disk.NewDiskStorage(strg, logger.Log, config.GetFileStoragePath())
	diskService := services.NewDiskService(diskStrg, config.GetStoreInterval(), config.GetRestore())
	diskService.FillMetricStorage()
	go diskService.CollectMetrics()

	r := routes.InitRouter(handler)

	logger.Log.Info("Running server", zap.String("address", config.Address))

	return http.ListenAndServe(config.Address, r)
}
