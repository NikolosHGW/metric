package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/services/metric"
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

	strg := memory.NewMemStorage()
	metricService := metric.NewMetricService(strg)
	handler := handlers.NewHandler(metricService)

	r := routes.InitRouter(handler)

	logger.Log.Info("Running server", zap.String("address", config.Address))

	return http.ListenAndServe(config.Address, r)
}
