package request

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_getStringValue(t *testing.T) {
	testCases := []struct {
		name     string
		v        interface{}
		expected string
	}{
		{"корректный случай #1", models.Gauge(3.14), "3.14"},
		{"корректный случай #2", models.Counter(42), "42"},
		{"некорректный тип", "foo", ""},
		{"нулевое значение", nil, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getStringValue(tc.v)
			if actual != tc.expected {
				t.Errorf("getStringValue(%v) = %q; want %q", tc.v, actual, tc.expected)
			}
		})
	}
}

func Test_getResultUrl(t *testing.T) {
	testCases := []struct {
		name        string
		hostAdrs    string
		metricType  string
		metricName  string
		metricValue string
		expected    string
	}{
		{"корректный случай #1", "localhost:8080", models.GaugeType, "Alloc", "123.45", "http://localhost:8080/update/gauge/Alloc/123.45"},
		{"корректный случай #2", "localhost:8080", models.CounterType, "PollCount", "42", "http://localhost:8080/update/counter/PollCount/42"},
		{"пустой тип метрики", "localhost:8080", "", "Alloc", "123.45", "http://localhost:8080/update/Alloc/123.45"},
		{"пустое имя метрики", "localhost:8080", models.GaugeType, "", "123.45", "http://localhost:8080/update/gauge/123.45"},
		{"пустое значение метрики", "localhost:8080", models.CounterType, "PollCounter", "", "http://localhost:8080/update/counter/PollCounter/"},
		{"пустые строки", "localhost:8080", "", "", "", "http://localhost:8080/update/"},
		{"пустой адрес", "", models.GaugeType, "Alloc", "123.45", "http:///update/gauge/Alloc/123.45"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getResultURL(tc.hostAdrs, tc.metricType, tc.metricName, tc.metricValue)
			if actual != tc.expected {
				t.Errorf("getResultUrl(%q, %q, %q, %q) = %q; want %q", tc.hostAdrs, tc.metricType, tc.metricName, tc.metricValue, actual, tc.expected)
			}
		})
	}
}

type MockClientMetrics struct {
	mock.Mock
}

func (m *MockClientMetrics) GetMetrics() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func TestGetMetricTypeMap(t *testing.T) {
	expectedMap := map[string]string{
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

	metricTypeMap := GetMetricTypeMap()

	assert.Equal(t, expectedMap, metricTypeMap, "The metric type map should match the expected values")
}

func TestSendMetrics(t *testing.T) {
	mockMetrics := new(MockClientMetrics)
	mockMetrics.On("GetMetrics").Return(map[string]interface{}{
		"test-metric": models.Gauge(123.456),
	})

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/test-metric/123.456", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		SendMetrics(ctx, mockMetrics, 1, testServer.URL[7:])
	}()

	<-ctx.Done()

	mockMetrics.AssertExpectations(t)

	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestSendJSONMetrics(t *testing.T) {
	mockMetrics := new(MockClientMetrics)
	expectedMetrics := map[string]interface{}{
		"test-metric": 123.456,
	}
	mockMetrics.On("GetMetrics").Return(expectedMetrics)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go SendJSONMetrics(ctx, mockMetrics, 1, "localhost", "test-key")

	<-ctx.Done()

	mockMetrics.AssertExpectations(t)
}

func TestSendBatchJSONMetrics(t *testing.T) {
	mockMetrics := new(MockClientMetrics)

	mockMetrics.On("GetMetrics").Return(map[string]interface{}{
		"testMetric": models.Gauge(42),
	})

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/updates/", req.URL.String())
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "gzip", req.Header.Get("Content-Encoding"))
		assert.Equal(t, "gzip", req.Header.Get("Accept-Encoding"))

		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)

		gr, err := gzip.NewReader(bytes.NewBuffer(body))
		assert.NoError(t, err)
		defer func() {
			err := gr.Close()
			if err != nil {
				t.Errorf("failed to close gzip reader: %v", err)
			}
		}()

		decompressedBody, err := io.ReadAll(gr)
		assert.NoError(t, err)

		var metrics []models.Metrics
		err = json.Unmarshal(decompressedBody, &metrics)
		assert.NoError(t, err)

		assert.Equal(t, "testMetric", metrics[0].ID)
		assert.Equal(t, models.GaugeType, metrics[0].MType)
		assert.NotNil(t, metrics[0].Value)
		assert.Equal(t, float64(42), *metrics[0].Value)

		_, err = rw.Write([]byte(`OK`))
		if err != nil {
			t.Errorf("failed to write body: %v", err)
		}
	}))
	defer server.Close()

	SendBatchJSONMetrics(mockMetrics, server.URL, "testKey")

	mockMetrics.AssertExpectations(t)
}
