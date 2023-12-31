package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/server/services/metric"
	"github.com/NikolosHGW/metric/internal/util"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type storageMock struct{}

func (sm storageMock) GetGaugeMetric(name string) (util.Gauge, error) {
	if name == "Alloc" {
		return 50.1, nil
	}

	return 0, errors.New("gauge metric not found")
}

func (sm storageMock) GetCounterMetric(name string) (util.Counter, error) {
	if name == "PollCounter" {
		return 50, nil
	}

	return 0, errors.New("counter metric not found")
}

func (sm storageMock) SetGaugeMetric(name string, value util.Gauge) {

}

func (sm storageMock) SetCounterMetric(name string, value util.Counter) {

}

func (sm storageMock) GetAllMetrics() []string {
	return []string{}
}

func TestWithSetMetricHandle(t *testing.T) {
	type want struct {
		code         int
		contentTypes []string
	}

	strg := storageMock{}
	metricService := metric.NewMetricService(strg)

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

			handler := NewHandler(metricService)
			handler.SetMetric(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			assert.ElementsMatch(t, test.want.contentTypes, res.Header.Values("Content-Type"))
		})
	}
}

func TestWithSetMetricHandle2(t *testing.T) {
	strg := storageMock{}
	metricService := metric.NewMetricService(strg)

	r := chi.NewRouter()

	handler := NewHandler(metricService)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.SetMetric)

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		code         int
		contentTypes []string
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
			defer res.Body.Close()

			assert.Equal(t, tc.want.code, res.StatusCode)

			assert.ElementsMatch(t, tc.want.contentTypes, res.Header.Values("Content-Type"))
		})
	}
}

func TestWithGetValueMetricHandle(t *testing.T) {
	strg := &storageMock{}
	metricService := metric.NewMetricService(strg)

	handler := NewHandler(metricService)

	r := chi.NewRouter()
	r.Get("/{metricType}/{metricName}", handler.GetValueMetric)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := ts.Client()

	tests := []struct {
		name       string
		metricType string
		metricName string
		wantStatus int
		wantBody   string
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
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, test.wantStatus, resp.StatusCode, "status code should match expected")

			assert.Equal(t, test.wantBody, string(body), "body should match expected")
		})
	}
}

func TestWithGetMetricsHandle(t *testing.T) {
	strg := storageMock{}
	metricService := metric.NewMetricService(strg)

	handler := NewHandler(metricService)

	r := chi.NewRouter()
	r.Get("/", handler.GetMetrics)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.NoError(t, err)

	client := ts.Client()
	resp, err := client.Do(req)

	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Contains(t, string(body), "<title>Список метрик</title>")
	assert.Contains(t, string(body), "<h1>Список метрик</h1>")
}
