package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NikolosHGW/metric/internal/models"
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

func SendJSONMetrics(m ClientMetrics, reportInterval int, host string) {
	metricTypeMap := util.GetMetricTypeMap()
	for {
		for k, v := range m.GetMetrics() {
			delta := getIntValue(metricTypeMap[k], v)
			value := getFloatValue(metricTypeMap[k], v)
			req := models.Metrics{
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

			resp, err := http.Post(getURL(host), "application/json", r)
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

func getURL(host string) string {
	sb := strings.Builder{}

	sb.WriteString("http://")
	sb.WriteString(host)
	sb.WriteString("/update")

	return sb.String()
}

func getIntValue(metricType string, value interface{}) int64 {
	if metricType == util.CounterType {
		v, ok := value.(util.Counter)
		if ok {
			return int64(v)
		}
	}

	return 0
}

func getFloatValue(metricType string, value interface{}) float64 {
	if metricType == util.GaugeType {
		v, ok := value.(util.Gauge)
		if ok {
			return float64(v)
		}
	}

	return 0
}
