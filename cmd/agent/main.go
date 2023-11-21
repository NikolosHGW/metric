package main

import (
	"time"

	"github.com/NikolosHGW/metric/internal/client/config"
	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/util"
)

func main() {
	config := config.NewConfig()

	parseFlags(&config.Flags)

	stats := metrics.NewMetrics()

	go util.CollectMetrics(stats, config.Flags.PollInterval)
	go util.SendMetrics(stats, config.Flags.ReportInterval, config.Flags.Endpoint.String())

	for {
		time.Sleep(10 * time.Second)
	}
}
