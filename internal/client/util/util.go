package util

import (
	"fmt"
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
)

const (
	pollInterval   = 2
	reportInterval = 10
)

func CollectMetrics(m metrics.ClientMetrics) {
	for {
		m.RefreshMetrics()
		m.IncPollCount()
		m.UpdateRandomValue()

		time.Sleep(pollInterval * time.Second)
	}
}

func SendMetrics(m metrics.ClientMetrics) {
	for {
		fmt.Println(m.GetMetrics())

		time.Sleep(reportInterval * time.Second)
	}
}
