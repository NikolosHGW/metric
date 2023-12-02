package main

import (
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/util"
)

func main() {
	config := NewConfig()

	stats := metrics.NewMetrics()

	go util.CollectMetrics(stats, config.GetPollInterval())
	go util.SendMetrics(stats, config.GetReportInterval(), config.GetAddress())

	for {
		time.Sleep(10 * time.Second)
	}
}
