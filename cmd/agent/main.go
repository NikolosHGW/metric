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

	reportInterval := config.GetReportInterval()
	adress := config.GetAddress()
	key := config.GetKey()

	for i := 0; i < rateLimit; i++ {
		go func() {
			for {
				requests <- struct{}{}
				request.SendBatchJSONMetrics(stats, adress, key)
				<-requests
				time.Sleep(time.Duration(reportInterval) * time.Second)
			}
		}()
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
