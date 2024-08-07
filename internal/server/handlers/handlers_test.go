package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/go-chi/chi"
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

func (m *mockLogger) Info(msg string, fields ...zap.Field) {}

func TestHandler_SetJSONMetric(t *testing.T) {
	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)

	handler := NewHandler(metricService, &mockLogger{})
	server := httptest.NewServer(http.HandlerFunc(handler.SetJSONMetric))
	defer server.Close()

	gaugeValue := f(42.1)
	counterValue := i(10)

	testCases := []struct {
		name     string
		request  models.Metrics
		expected models.Metrics
		status   int
	}{
		{
			name:     "положительный тест: установить gauge метрику",
			request:  models.Metrics{ID: "foo", MType: models.GaugeType, Value: gaugeValue},
			expected: models.Metrics{ID: "foo", MType: models.GaugeType, Value: gaugeValue},
			status:   http.StatusOK,
		},
		{
			name:     "положительный тест: установить counter метрику",
			request:  models.Metrics{ID: "bar", MType: models.CounterType, Delta: counterValue},
			expected: models.Metrics{ID: "bar", MType: models.CounterType, Delta: counterValue},
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
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Errorf("failed body close: %v", err)
				}
			}()

			assert.Equal(t, tc.status, resp.StatusCode)

			var actual models.Metrics
			err = json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestHandler_GetMetric(t *testing.T) {
	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)
	err := metricService.SetJSONMetric(context.Background(), models.Metrics{ID: "cpu", MType: "gauge", Value: f(0.5)})
	if err != nil {
		t.Errorf("failed SetJSONMetric: %v", err)
	}
	err = metricService.SetJSONMetric(context.Background(), models.Metrics{ID: "memory", MType: "counter", Delta: i(10)})
	if err != nil {
		t.Errorf("2 failed SetJSONMetric: %v", err)
	}
	h := NewHandler(metricService, &mockLogger{})

	tests := []struct {
		name           string
		requestBody    string
		expectedHeader string
		expectedBody   string
		expectedStatus int
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
			expectedBody:   `{"value":0.5,"id":"cpu","type":"gauge"}`,
		},
		{
			name:           "положительный тест: получение существующей метрики counter",
			requestBody:    `{"id": "memory", "type": "counter"}`,
			expectedStatus: http.StatusOK,
			expectedHeader: "application/json",
			expectedBody:   `{"delta":10,"id":"memory","type":"counter"}`,
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

type storageMock struct{}

func (sm storageMock) GetGaugeMetric(_ context.Context, name string) (models.Gauge, error) {
	if name == "Alloc" {
		return 50.1, nil
	}

	return 0, errors.New("gauge metric not found")
}

func (sm storageMock) GetCounterMetric(_ context.Context, name string) (models.Counter, error) {
	if name == "PollCounter" {
		return 50, nil
	}

	return 0, errors.New("counter metric not found")
}

func (sm storageMock) SetGaugeMetric(_ context.Context, name string, value models.Gauge) error {
	return nil
}

func (sm storageMock) SetCounterMetric(_ctx context.Context, name string, value models.Counter) error {
	return nil
}

func (sm storageMock) GetAllMetrics(_ctx context.Context) []string {
	return []string{}
}

func (sm storageMock) SetMetric(_ctx context.Context, m models.Metrics) error { return nil }

func (sm storageMock) GetMetric(_ctx context.Context, name string) (models.Metrics, error) {
	return models.Metrics{}, nil
}

func (sm storageMock) GetIsDBConnected() bool {
	return false
}

func (sm storageMock) UpsertMetrics(ctx context.Context, mc models.MetricCollection) (models.MetricCollection, error) {
	return *models.NewMetricCollection(), nil
}

func TestWithSetMetricHandle(t *testing.T) {
	type want struct {
		contentTypes []string
		code         int
	}

	strg := storageMock{}
	metricService := services.NewMetricService(strg)

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "положительный тест для метрики типа counter",
			url:  "/update/counter/someMetric/527",
			want: want{
				code:         200,
				contentTypes: []string{"text/plain", "charset=utf-8"},
			},
		},
		{
			name: "положительный тест для метрики типа gauge",
			url:  "/update/gauge/someMetric/527.2",
			want: want{
				code:         200,
				contentTypes: []string{"text/plain", "charset=utf-8"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()

			handler := NewHandler(metricService, &mockLogger{})
			handler.SetMetric(w, request)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {
					t.Errorf("failed body close: %v", err)
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)

			assert.ElementsMatch(t, test.want.contentTypes, res.Header.Values("Content-Type"))
		})
	}
}

func TestWithSetMetricHandle2(t *testing.T) {
	strg := storageMock{}
	metricService := services.NewMetricService(strg)

	r := chi.NewRouter()

	handler := NewHandler(metricService, &mockLogger{})

	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.SetMetric)

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		contentTypes []string
		code         int
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "положительный тест для метрики типа counter",
			url:  "/update/counter/someMetric/527",
			want: want{
				code:         200,
				contentTypes: []string{"text/plain", "charset=utf-8"},
			},
		},
		{
			name: "положительный тест для метрики типа gauge",
			url:  "/update/gauge/someMetric/527.2",
			want: want{
				code:         200,
				contentTypes: []string{"text/plain", "charset=utf-8"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, ts.URL+tc.url, nil)
			assert.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer func() {
				err := res.Body.Close()
				if err != nil {
					t.Errorf("failed body close: %v", err)
				}
			}()

			assert.Equal(t, tc.want.code, res.StatusCode)

			assert.ElementsMatch(t, tc.want.contentTypes, res.Header.Values("Content-Type"))
		})
	}
}

func TestWithGetValueMetricHandle(t *testing.T) {
	strg := &storageMock{}
	metricService := services.NewMetricService(strg)

	handler := NewHandler(metricService, &mockLogger{})

	r := chi.NewRouter()
	r.Get("/{metricType}/{metricName}", handler.GetValueMetric)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := ts.Client()

	tests := []struct {
		name       string
		metricType string
		metricName string
		wantBody   string
		wantStatus int
	}{
		{
			name:       "положительный тест для метрики типа gauge",
			metricType: "gauge",
			metricName: "Alloc",
			wantStatus: http.StatusOK,
			wantBody:   "50.1",
		},
		{
			name:       "положительный тест для метрики типа counter",
			metricType: "counter",
			metricName: "PollCounter",
			wantStatus: http.StatusOK,
			wantBody:   "50",
		},
		{
			name:       "отрицательный тест с несуществующим типом метрики",
			metricType: "string",
			metricName: "Alloc",
			wantStatus: http.StatusNotFound,
			wantBody:   "",
		},
		{
			name:       "отрицательный тест с несуществующим именем метрики",
			metricType: "gauge",
			metricName: "qwerty",
			wantStatus: http.StatusNotFound,
			wantBody:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+test.metricType+"/"+test.metricName, nil)
			assert.NoError(t, err)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Errorf("failed body close: %v", err)
				}
			}()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, test.wantStatus, resp.StatusCode, "status code should match expected")

			assert.Equal(t, test.wantBody, string(body), "body should match expected")
		})
	}
}

func TestWithGetMetricsHandle(t *testing.T) {
	strg := storageMock{}
	metricService := services.NewMetricService(strg)

	handler := NewHandler(metricService, &mockLogger{})

	r := chi.NewRouter()
	r.Get("/", handler.GetMetrics)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.NoError(t, err)

	client := ts.Client()
	resp, err := client.Do(req)

	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Errorf("failed body close: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Contains(t, string(body), "<title>Список метрик</title>")
	assert.Contains(t, string(body), "<h1>Список метрик</h1>")
}

func TestHandler_UpsertMetrics(t *testing.T) {
	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)

	handler := NewHandler(metricService, &mockLogger{})
	server := httptest.NewServer(http.HandlerFunc(handler.UpsertMetrics))
	defer server.Close()

	gaugeValue := f(42.1)
	counterValue := i(10)

	testCases := []struct {
		name     string
		request  []models.Metrics
		expected []models.Metrics
		status   int
	}{
		{
			name: "положительный тест: добавить пачку метрик",
			request: []models.Metrics{
				{ID: "foo", MType: models.GaugeType, Value: gaugeValue},
				{ID: "bar", MType: models.CounterType, Delta: counterValue},
			},
			expected: []models.Metrics{
				{ID: "foo", MType: models.GaugeType, Value: gaugeValue},
				{ID: "bar", MType: models.CounterType, Delta: counterValue},
			},
			status: http.StatusOK,
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
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Errorf("failed body close: %v", err)
				}
			}()

			assert.Equal(t, tc.status, resp.StatusCode)

			var actual []models.Metrics
			err = json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
