package main

import (
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/request"
)

func main() {
	config := NewConfig()

	stats := metrics.NewMetrics()

	pollTicker := time.NewTicker(time.Duration(config.GetPollInterval()) * time.Second)

	go func() {
		for range pollTicker.C {
			stats.CollectMetrics()
			stats.CollectAdvancedMetric()
		}
	}()

	rateLimit := config.GetRateLimit()

	requests := make(chan struct{}, rateLimit)

	reportTicker := time.NewTicker(time.Duration(config.GetReportInterval()) * time.Second)

	for i := 0; i < rateLimit; i++ {
		go func() {
			for range reportTicker.C {
				requests <- struct{}{}
				request.SendBatchJSONMetrics(stats, config.GetAddress(), config.GetKey())
				<-requests
			}
		}()
	}

	select {}
}
