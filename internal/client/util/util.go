package util

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/util"
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
	metricTypeMap := util.GetMetricTypeMap()
	for {
		for k, v := range m.GetMetrics() {
			result := getStringValue(v)
			adrs := getResultUrl(metricTypeMap[k], k, result)

			resp, err := http.Post(adrs, "text/plain", nil)
			if err != nil {
				continue
			}
			resp.Body.Close()
		}

		time.Sleep(reportInterval * time.Second)
	}
}

func getStringValue(v interface{}) string {
	switch v2 := v.(type) {
	case util.Gauge:
		return strconv.FormatFloat(float64(v2), 'f', -1, 64)
	case util.Counter:
		return strconv.Itoa(int(v2))
	}

	return ""
}

func getResultUrl(metricType string, metricName string, metricValue string) string {
	sb := strings.Builder{}

	sb.WriteString("http://localhost:8080/update/")
	sb.WriteString(metricType)
	sb.WriteString("/")
	sb.WriteString(metricName)
	sb.WriteString("/")
	sb.WriteString(metricValue)

	return sb.String()
}