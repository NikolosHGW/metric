package main

import (
	"fmt"
	"time"

	"github.com/NikolosHGW/metric/internal/client/config"
	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/request"
)

const defaultTagValue = "N/A"

var (
	buildVersion = defaultTagValue
	buildDate    = defaultTagValue
	buildCommit  = defaultTagValue
)

func main() {
	config := config.NewConfig()

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
				request.SendBatchJSONMetrics(stats, config.GetAddress(), config.GetKey(), config.GetCryptoKeyPath())
				<-requests
			}
		}()
	}

	fmt.Println(
		"Build version: ", buildVersion, "\n",
		"Build date: ", buildDate, "\n",
		"Build commit: ", buildCommit,
	)

	select {}
}
