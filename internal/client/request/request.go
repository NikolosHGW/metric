package request

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NikolosHGW/metric/internal/models"
)

type ClientMetrics interface {
	GetMetrics() map[string]interface{}
}

func GetMetricTypeMap() map[string]string {
	return map[string]string{
		models.Alloc:           models.GaugeType,
		models.BuckHashSys:     models.GaugeType,
		models.Frees:           models.GaugeType,
		models.GCCPUFraction:   models.GaugeType,
		models.GCSys:           models.GaugeType,
		models.HeapAlloc:       models.GaugeType,
		models.HeapIdle:        models.GaugeType,
		models.HeapInuse:       models.GaugeType,
		models.HeapObjects:     models.GaugeType,
		models.HeapReleased:    models.GaugeType,
		models.HeapSys:         models.GaugeType,
		models.LastGC:          models.GaugeType,
		models.Lookups:         models.GaugeType,
		models.MCacheInuse:     models.GaugeType,
		models.MCacheSys:       models.GaugeType,
		models.MSpanInuse:      models.GaugeType,
		models.MSpanSys:        models.GaugeType,
		models.Mallocs:         models.GaugeType,
		models.NextGC:          models.GaugeType,
		models.NumForcedGC:     models.GaugeType,
		models.NumGC:           models.GaugeType,
		models.OtherSys:        models.GaugeType,
		models.PauseTotalNs:    models.GaugeType,
		models.StackInuse:      models.GaugeType,
		models.StackSys:        models.GaugeType,
		models.Sys:             models.GaugeType,
		models.TotalAlloc:      models.GaugeType,
		models.PollCount:       models.CounterType,
		models.RandomValue:     models.GaugeType,
		models.TotalMemory:     models.GaugeType,
		models.FreeMemory:      models.GaugeType,
		models.CPUutilization1: models.GaugeType,
	}
}

func SendMetrics(ctx context.Context, m ClientMetrics, reportInterval int, host string) {
	metricTypeMap := GetMetricTypeMap()
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for k, v := range m.GetMetrics() {
				result := getStringValue(v)
				adrs := getResultURL(host, metricTypeMap[k], k, result)

				resp, err := http.Post(adrs, "text/plain", nil)
				if err != nil {
					continue
				}
				err = resp.Body.Close()
				if err != nil {
					log.Println("can not close body SendMetrics", err)
					continue
				}
			}
		}
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

func SendJSONMetrics(ctx context.Context, m ClientMetrics, reportInterval int, host, key string) {
	metricTypeMap := GetMetricTypeMap()
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for k, v := range m.GetMetrics() {
				delta := GetIntValue(metricTypeMap[k], v)
				value := GetFloatValue(metricTypeMap[k], v)
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

				hash := hash(data, key)

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
				if hash != "" {
					nr.Header.Set("HashSHA256", hash)
				}
				resp, err := http.DefaultClient.Do(nr)

				if err != nil {
					log.Println("metric/internal/client/util/util.go SendMetrics cannot Post", err)
					continue
				}
				log.Println("metric/internal/client/util/util.go SendMetrics post status", resp.Status)
				err = resp.Body.Close()
				if err != nil {
					log.Println("can not close body SendJSONMetrics", err)
					continue
				}
			}
		}
	}
}

func SendBatchJSONMetrics(m ClientMetrics, host, key string) {
	metricTypeMap := GetMetricTypeMap()
	metricsBatch := make([]models.Metrics, 0, len(m.GetMetrics()))

	for k, v := range m.GetMetrics() {
		delta := GetIntValue(metricTypeMap[k], v)
		value := GetFloatValue(metricTypeMap[k], v)
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
		return
	}

	hash := hash(data, key)

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, gzipErr := zb.Write(data)
	if gzipErr != nil {
		log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot gzip write", gzipErr)
		return
	}
	err = zb.Close()
	if err != nil {
		log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot gzip close", err)
		return
	}

	nr, err := http.NewRequest(http.MethodPost, getUpdatesURL(host), buf)
	if err != nil {
		log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot create NewRequest", err)
		return
	}
	nr.Header.Set("Content-Type", "application/json")
	nr.Header.Set("Content-Encoding", "gzip")
	nr.Header.Set("Accept-Encoding", "gzip")
	if hash != "" {
		nr.Header.Set("HashSHA256", hash)
	}
	resp, err := http.DefaultClient.Do(nr)

	if err != nil {
		log.Println("metric/internal/client/util/util.go SendBatchMetrics cannot Post", err)
		return
	}
	log.Println("metric/internal/client/util/util.go SendBatchMetrics post status", resp.Status)
	err = resp.Body.Close()
	if err != nil {
		log.Println("can not close body SendBatchJSONMetrics", err)
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
	sb.WriteString("/updates/")

	return sb.String()
}

func GetIntValue(metricType string, value interface{}) int64 {
	if metricType == models.CounterType {
		v, ok := value.(models.Counter)
		if ok {
			return int64(v)
		}
	}

	return 0
}

func GetFloatValue(metricType string, value interface{}) float64 {
	if metricType == models.GaugeType {
		v, ok := value.(models.Gauge)
		if ok {
			return float64(v)
		}
	}

	return 0
}

func hash(data []byte, key string) string {
	if key != "" {
		h := hmac.New(sha256.New, []byte(key))
		h.Write(data)
		return hex.EncodeToString(h.Sum(nil))
	}

	return ""
}
