package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type MockMetricService struct{}

func (m *MockMetricService) SetMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	return nil
}

func (m *MockMetricService) SetJSONMetric(ctx context.Context, metric models.Metrics) error {
	return nil
}

func (m *MockMetricService) GetMetricValue(ctx context.Context, metricType, metricName string) (string, error) {
	return "123", nil
}

func (m *MockMetricService) GetMetricByName(ctx context.Context, metricName string) (models.Metrics, error) {
	return models.Metrics{ID: metricName, MType: "gauge", Value: f(123)}, nil
}

func (m *MockMetricService) GetAllMetrics(ctx context.Context) []string {
	return []string{"metric1", "metric2"}
}

func (m *MockMetricService) GetIsDBConnected() bool {
	return true
}

func (m *MockMetricService) UpsertMetrics(
	ctx context.Context,
	metricCollection models.MetricCollection,
) (models.MetricCollection, error) {
	return metricCollection, nil
}

type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {}

func ExampleHandler_SetMetric() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	req, _ := http.NewRequest("POST", "/metrics/gauge/metricName/123", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Post("/metrics/{metricType}/{metricName}/{metricValue}", handler.SetMetric)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to set metric")
	}

	fmt.Println(rr.Code)

	// Output:
	// 200
}

func ExampleHandler_GetValueMetric() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	req, _ := http.NewRequest("GET", "/metrics/gauge/metricName", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/metrics/{metricType}/{metricName}", handler.GetValueMetric)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to get metric value")
	}

	fmt.Println(rr.Body.String())

	// Output:
	// 123
}

func ExampleHandler_SetJSONMetric() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	metric := models.Metrics{ID: "metricName", MType: "gauge", Value: f(123)}
	jsonMetric, _ := json.Marshal(metric)
	req, _ := http.NewRequest("POST", "/update", bytes.NewBuffer(jsonMetric))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Post("/update", handler.SetJSONMetric)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to set JSON metric")
	}

	fmt.Println(rr.Code)

	// Output:
	// 200
}

func ExampleHandler_GetMetric() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	metric := models.Metrics{ID: "metricName", MType: "gauge"}
	jsonMetric, _ := json.Marshal(metric)
	req, _ := http.NewRequest("POST", "/value", bytes.NewBuffer(jsonMetric))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Post("/value", handler.GetMetric)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		fmt.Println("failed to get metric, status code:", rr.Code)
		panic("failed to get metric")
	}

	fmt.Println(rr.Body.String())

	// Output:
	// {"value":123,"id":"metricName","type":"gauge"}
}

func ExampleHandler_GetMetrics() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/", handler.GetMetrics)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to get metrics")
	}

	fmt.Println(rr.Header().Get("Content-Type"))

	// Output:
	// text/html
}

func ExampleHandler_PingDB() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	req, _ := http.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/ping", handler.PingDB)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to ping DB")
	}

	fmt.Println(rr.Code)

	// Output:
	// 200
}

func ExampleHandler_UpsertMetrics() {
	ms := &MockMetricService{}
	logger := &MockLogger{}
	handler := NewHandler(ms, logger)

	metrics := []models.Metrics{
		{ID: "metric1", MType: "gauge", Value: f(123)},
		{ID: "metric2", MType: "counter", Delta: i(10)},
	}
	jsonMetrics, _ := json.Marshal(metrics)
	req, _ := http.NewRequest("POST", "/updates", bytes.NewBuffer(jsonMetrics))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Post("/updates", handler.UpsertMetrics)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to upsert metrics")
	}

	fmt.Println(rr.Code)

	// Output:
	// 200
}
