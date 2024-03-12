package main

import (
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/request"
)

func main() {
	config := NewConfig()

	stats := metrics.NewMetrics()

	go request.CollectMetrics(stats, config.GetPollInterval())
	go request.SendJSONMetrics(stats, config.GetReportInterval(), config.GetAddress())
	go request.SendBatchJSONMetrics(stats, config.GetReportInterval(), config.GetAddress())

	for {
		time.Sleep(10 * time.Second)
	}
}
