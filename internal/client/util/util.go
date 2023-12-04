package util

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
			result := getStringValue(v)
			adrs := getResultURL(host, metricTypeMap[k], k, result)

			resp, err := http.Post(adrs, "text/plain", nil)
			if err != nil {
				continue
			}
			resp.Body.Close()
		}

		time.Sleep(time.Duration(reportInterval) * time.Second)
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

func getResultURL(host, metricType, metricName, metricValue string) string {
	sb := strings.Builder{}

	sb.WriteString("http://")
	sb.WriteString(host)
	sb.WriteString("/update/")
	if metricType != "" {
		sb.WriteString(metricType)
		sb.WriteString("/")
	}
	if metricName != "" {
		sb.WriteString(metricName)
		sb.WriteString("/")
	}
	if metricValue != "" {
		sb.WriteString(metricValue)
	}

	return sb.String()
}
