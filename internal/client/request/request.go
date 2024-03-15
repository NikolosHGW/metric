package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NikolosHGW/metric/internal/models"
)

type ClientMetrics interface {
	IncPollCount()
	UpdateRandomValue()
	RefreshMetrics()
	GetMetrics() map[string]interface{}
}

func GetMetricTypeMap() map[string]string {
	return map[string]string{
		models.Alloc:         models.GaugeType,
		models.BuckHashSys:   models.GaugeType,
		models.Frees:         models.GaugeType,
		models.GCCPUFraction: models.GaugeType,
		models.GCSys:         models.GaugeType,
		models.HeapAlloc:     models.GaugeType,
		models.HeapIdle:      models.GaugeType,
		models.HeapInuse:     models.GaugeType,
		models.HeapObjects:   models.GaugeType,
		models.HeapReleased:  models.GaugeType,
		models.HeapSys:       models.GaugeType,
		models.LastGC:        models.GaugeType,
		models.Lookups:       models.GaugeType,
		models.MCacheInuse:   models.GaugeType,
		models.MCacheSys:     models.GaugeType,
		models.MSpanInuse:    models.GaugeType,
		models.MSpanSys:      models.GaugeType,
		models.Mallocs:       models.GaugeType,
		models.NextGC:        models.GaugeType,
		models.NumForcedGC:   models.GaugeType,
		models.NumGC:         models.GaugeType,
		models.OtherSys:      models.GaugeType,
		models.PauseTotalNs:  models.GaugeType,
		models.StackInuse:    models.GaugeType,
		models.StackSys:      models.GaugeType,
		models.Sys:           models.GaugeType,
		models.TotalAlloc:    models.GaugeType,
		models.PollCount:     models.CounterType,
		models.RandomValue:   models.GaugeType,
	}
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
	metricTypeMap := GetMetricTypeMap()
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
	case models.Gauge:
		return strconv.FormatFloat(float64(v2), 'f', -1, 64)
	case models.Counter:
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
	metricTypeMap := GetMetricTypeMap()
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

			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, gzipErr := zb.Write(data)
			if gzipErr != nil {
				log.Println("metric/internal/client/util/util.go SendMetrics cannot gzip write", gzipErr)
				continue
			}
			err = zb.Close()
			if err != nil {
				log.Println("metric/internal/client/util/util.go SendMetrics cannot gzip cloze", err)
				continue
			}

			nr, err := http.NewRequest(http.MethodPost, getURL(host), buf)
			if err != nil {
				log.Println("metric/internal/client/util/util.go SendMetrics cannot create NewRequest", err)
				continue
			}
			nr.Header.Set("Content-Type", "application/json")
			nr.Header.Set("Content-Encoding", "gzip")
			nr.Header.Set("Accept-Encoding", "gzip")
			resp, err := http.DefaultClient.Do(nr)

			if err != nil {
				log.Println("metric/internal/client/util/util.go SendMetrics cannot Post", err)
				continue
			}
			log.Println("metric/internal/client/util/util.go SendMetrics: ", data, "; ", "post status", resp.Status)
			resp.Body.Close()
		}

		time.Sleep(time.Duration(reportInterval) * time.Second)
	}
}

func SendBatchJSONMetrics(m ClientMetrics, reportInterval int, host string) {
	metricTypeMap := GetMetricTypeMap()
	for {
		var metricsBatch []models.Metrics

		for k, v := range m.GetMetrics() {
			delta := getIntValue(metricTypeMap[k], v)
			value := getFloatValue(metricTypeMap[k], v)
			metric := models.Metrics{
				ID:    k,
				MType: metricTypeMap[k],
				Delta: &delta,
				Value: &value,
			}
			metricsBatch = append(metricsBatch, metric)
		}

		data, err := json.Marshal(metricsBatch)
		if err != nil {
			log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot Marshal", err)
			continue
		}

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, gzipErr := zb.Write(data)
		if gzipErr != nil {
			log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot gzip write", gzipErr)
			continue
		}
		err = zb.Close()
		if err != nil {
			log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot gzip close", err)
			continue
		}

		nr, err := http.NewRequest(http.MethodPost, getUpdatesURL(host), buf)
		if err != nil {
			log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot create NewRequest", err)
			continue
		}
		nr.Header.Set("Content-Type", "application/json")
		nr.Header.Set("Content-Encoding", "gzip")
		nr.Header.Set("Accept-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(nr)

		if err != nil {
			log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot Post", err)
			continue
		}
		log.Println("metric/internal/client/util/util.go SendBatchMetrics: ", data, "; ", "post status", resp.Status)
		resp.Body.Close()

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

func getUpdatesURL(host string) string {
	sb := strings.Builder{}

	sb.WriteString("http://")
	sb.WriteString(host)
	sb.WriteString("/updates")

	return sb.String()
}

func getIntValue(metricType string, value interface{}) int64 {
	if metricType == models.CounterType {
		v, ok := value.(models.Counter)
		if ok {
			return int64(v)
		}
	}

	return 0
}

func getFloatValue(metricType string, value interface{}) float64 {
	if metricType == models.GaugeType {
		v, ok := value.(models.Gauge)
		if ok {
			return float64(v)
		}
	}

	return 0
}
