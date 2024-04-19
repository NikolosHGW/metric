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

	rateLimit := config.GetRateLimit()

	requests := make(chan struct{}, rateLimit)

	for i := 0; i < rateLimit; i++ {
		go func() {
			for {
				requests <- struct{}{}
				request.SendBatchJSONMetrics(stats, config.GetReportInterval(), config.GetAddress(), config.GetKey())
				<-requests
				time.Sleep(time.Duration(config.GetReportInterval()) * time.Second)
			}
		}()
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
