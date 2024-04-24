package main

import (
	"log"
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/request"
	"github.com/NikolosHGW/metric/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	config := NewConfig()

	stats := metrics.NewMetrics()

	go request.CollectMetrics(stats, config.GetPollInterval())
	go func() {
		for {
			stats.LockMutex()
			v, err := mem.VirtualMemory()
			if err != nil {
				log.Println("failed virtual memory metrics", err)
			}
			cpuPercentages, err := cpu.Percent(0, false)
			if err != nil {
				log.Println("failed cpu percent", err)
			}
			stats.SetAdvanceMetrics(models.Gauge(v.Total), models.Gauge(v.Free), models.Gauge(cpuPercentages[0]))
			time.Sleep(time.Duration(config.GetPollInterval()) * time.Second)
			stats.UnlockMutex()
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

	go request.SendJSONMetrics(stats, reportInterval, adress, key)

	for {
		time.Sleep(10 * time.Second)
	}
}
