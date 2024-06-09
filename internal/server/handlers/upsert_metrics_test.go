package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/client/request"
	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
)

func BenchmarkUpsertMetrics(b *testing.B) {
	metricCollectionBuffer := createTestMetricCollectionJSON(b)
	req := httptest.NewRequest(http.MethodPost, "/updates/", metricCollectionBuffer)
	w := httptest.NewRecorder()

	h := getHandler()

	for i := 0; i < b.N; i++ {
		h.UpsertMetrics(w, req)
	}
}

func createTestMetricCollectionJSON(b *testing.B) *bytes.Buffer {
	stats := metrics.NewMetrics()
	stats.CollectMetrics()
	stats.CollectAdvancedMetric()

	metricTypeMap := request.GetMetricTypeMap()
	var metricsBatch []models.Metrics

	for k, v := range stats.GetMetrics() {
		delta := request.GetIntValue(metricTypeMap[k], v)
		value := request.GetFloatValue(metricTypeMap[k], v)
		metric := models.Metrics{
			ID:    k,
			MType: metricTypeMap[k],
			Delta: &delta,
			Value: &value,
		}
		metricsBatch = append(metricsBatch, metric)
	}

	body, err := json.Marshal(metricsBatch)
	if err != nil {
		b.Fatal(err)
	}

	return bytes.NewBuffer(body)
}

func getHandler() *Handler {
	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)

	return NewHandler(metricService, &mockLogger{})
}
