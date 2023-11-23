package main

import (
	"time"

	"github.com/NikolosHGW/metric/internal/client/config"
	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/util"
)

func main() {
	config := config.NewConfig()

	parseFlags(config)
	config.InitEnv()

	stats := metrics.NewMetrics()

	go util.CollectMetrics(stats, config.PollInterval)
	go util.SendMetrics(stats, config.ReportInterval, config.Address)

	for {
		time.Sleep(10 * time.Second)
	}
}
