package handlers

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/server/services/metric"
	"github.com/NikolosHGW/metric/internal/server/storage/memory"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func f(v float64) *float64 {
	return &v
}

func i(v int64) *int64 {
	return &v
}

type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields ...zap.Field) {}

func TestHandler_SetMetric(t *testing.T) {
	strg := memory.NewMemStorage()
	metricService := metric.NewMetricService(strg)

	handler := NewHandler(metricService, &mockLogger{})
	server := httptest.NewServer(http.HandlerFunc(handler.SetMetric))
	defer server.Close()

	testCases := []struct {
		name     string
		request  models.Metrics
		expected models.Metrics
		status   int
	}{
		{
			name:     "положительный тест: установить gauge метрику",
			request:  models.Metrics{ID: "foo", MType: "gauge", Value: f(42.1)},
			expected: models.Metrics{ID: "foo", MType: "gauge", Value: f(42.1)},
			status:   http.StatusOK,
		},
		{
			name:     "положительный тест: установить counter метрику",
			request:  models.Metrics{ID: "bar", MType: "counter", Delta: i(10)},
			expected: models.Metrics{ID: "bar", MType: "counter", Delta: i(10)},
			status:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tc.request)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(reqBody))
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.status, resp.StatusCode)

			var actual models.Metrics
			err = json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestHandler_GetMetric(t *testing.T) {
	strg := memory.NewMemStorage()
	metricService := metric.NewMetricService(strg)
	metricService.SetMetric(models.Metrics{ID: "cpu", MType: "gauge", Value: f(0.5)})
	metricService.SetMetric(models.Metrics{ID: "memory", MType: "counter", Delta: i(10)})
	h := NewHandler(metricService, &mockLogger{})

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "отрицательный тест: невалидный JSON",
			requestBody:    `{"id": "cpu", "type": "gauge", value: 0.5}`, // нет кавычек у ключа "value"
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "text/plain; charset=utf-8",
			expectedBody:   "неверный формат запроса\n",
		},
		{
			name:           "отрицательный тест: невалидный тип метрики",
			requestBody:    `{"id": "cpu", "type": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "text/plain; charset=utf-8",
			expectedBody:   "неверный формат запроса\n",
		},
		{
			name:           "отрицательный тест: метрика не найдена",
			requestBody:    `{"id": "disk", "type": "gauge"}`,
			expectedStatus: http.StatusNotFound,
			expectedHeader: "text/plain; charset=utf-8",
			expectedBody:   "метрика не найдена\n",
		},
		{
			name:           "положительный тест: получение существующей метрики gauge",
			requestBody:    `{"id": "cpu", "type": "gauge"}`,
			expectedStatus: http.StatusOK,
			expectedHeader: "application/json",
			expectedBody:   `{"id":"cpu","type":"gauge","value":0.5}`,
		},
		{
			name:           "положительный тест: получение существующей метрики counter",
			requestBody:    `{"id": "memory", "type": "counter"}`,
			expectedStatus: http.StatusOK,
			expectedHeader: "application/json",
			expectedBody:   `{"id":"memory","type":"counter","delta":10}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/value", bytes.NewBufferString(test.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			h.GetMetric(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)

			assert.Contains(t, rr.Header().Get("Content-Type"), test.expectedHeader)

			assert.Equal(t, test.expectedBody, rr.Body.String())
		})
	}
}
