package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/services/metric"
	"github.com/NikolosHGW/metric/internal/server/storage/memory"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(fmt.Errorf("server/main run_ListenAndServe: %w", err))
	}
}

func run() error {
	config := NewConfig()

	strg := memory.NewMemStorage()
	metricService := metric.NewMetricService(strg)
	handler := handlers.NewHandler(metricService)

	r := routes.InitRouter(handler)

	return http.ListenAndServe(config.Address, r)
}
