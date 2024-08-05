package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
	defer pollTicker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-pollTicker.C:
				stats.CollectMetrics()
				stats.CollectAdvancedMetric()
			case <-ctx.Done():
				return
			}
		}
	}()

	rateLimit := config.GetRateLimit()

	requests := make(chan struct{}, rateLimit)

	reportTicker := time.NewTicker(time.Duration(config.GetReportInterval()) * time.Second)
	defer reportTicker.Stop()

	for i := 0; i < rateLimit; i++ {
		go func() {
			for {
				select {
				case <-reportTicker.C:
					requests <- struct{}{}
					request.SendBatchJSONMetrics(stats, config.GetAddress(), config.GetKey(), config.GetCryptoKeyPath())
					<-requests
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	fmt.Println(
		"Build version: ", buildVersion, "\n",
		"Build date: ", buildDate, "\n",
		"Build commit: ", buildCommit,
	)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sig := <-signalChan
	fmt.Println("Received signal:", sig)
	cancel()

	time.Sleep(2 * time.Second)

	fmt.Println("Agent exited gracefully")
}
