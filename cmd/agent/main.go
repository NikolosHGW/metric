package main

import (
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/util"
)

func main() {
	stats := metrics.NewMetrics()

	go util.CollectMetrics(stats)
	go util.SendMetrics(stats)

	for {
		time.Sleep(10 * time.Second)
	}
}
