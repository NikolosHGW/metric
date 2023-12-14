package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/NikolosHGW/metric/internal/models/metric"
	"github.com/NikolosHGW/metric/internal/util"
)

type ClientMetrics interface {
	IncPollCount()
	UpdateRandomValue()
	RefreshMetrics()
	GetMetrics() map[string]interface{}
}

func CollectMetrics(m ClientMetrics, pollInterval int) {
	for {
		m.RefreshMetrics()
		m.IncPollCount()
		m.UpdateRandomValue()

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func SendMetrics(m ClientMetrics, reportInterval int, host string) {
	metricTypeMap := util.GetMetricTypeMap()
	for {
		for k, v := range m.GetMetrics() {
			delta := getIntValue(metricTypeMap[k], v)
			value := getFloatValue(metricTypeMap[k], v)
			req := metric.Metrics{
				ID:    k,
				MType: metricTypeMap[k],
				Delta: &delta,
				Value: &value,
			}

			data, err := json.Marshal(req)
			if err != nil {
				log.Println("metric/internal/client/util/util.go SendMetrics cannot Marshal", err)
				continue
			}
			r := bytes.NewReader(data)

			resp, err := http.Post("/update", "application/json", r)
			if err != nil {
				log.Println("metric/internal/client/util/util.go SendMetrics cannot Post", err)
				continue
			}
			log.Println("metric/internal/client/util/util.go SendMetrics post status", resp.Status)
			resp.Body.Close()
		}

		time.Sleep(time.Duration(reportInterval) * time.Second)
	}
}

func getIntValue(metricType string, value interface{}) int64 {
	if metricType == util.CounterType {
		v, ok := value.(int64)
		if ok {
			return v
		}
	}

	return 0
}

func getFloatValue(metricType string, value interface{}) float64 {
	if metricType == util.GaugeType {
		v, ok := value.(float64)
		if ok {
			return v
		}
	}

	return 0
}
