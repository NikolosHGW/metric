package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/config"
	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/storage/memory"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(fmt.Errorf("server/main run_ListenAndServe: %w", err))
	}
}

func run() error {
	config := config.NewConfig()

	parseFlags(&config.Endpoint)

	strg := memory.NewMemStorage()

	r := routes.InitRouter(strg)

	return http.ListenAndServe(config.Endpoint.String(), r)
}
