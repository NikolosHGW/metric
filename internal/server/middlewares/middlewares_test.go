package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getFakeHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestCheckMetricNameMiddleware(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "положительный тест: эндпоинт с типом метрики, с именем метрики и со значением",
			url:  "/update/counter/someMetric/527",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "отрицительный тест: эндпоинт без имя метрики",
			url:  "/update/counter/527",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()

			CheckMetricNameMiddleware(getFakeHandler()).ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestCheckTypeAndValueMiddleware(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "положительный тест: тип метрики counter и тип значения int64 (кастомный Counter)",
			url:  "/update/counter/someMetric/527",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "положительный тест: тип метрики gauge и тип значения float64 (кастомный Gauge)",
			url:  "/update/gauge/someMetric/527.2",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "отрицительный тест: тип метрики someType (ни counter, ни gauge)",
			url:  "/update/someType/someMetric/527",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "отрицительный тест: тип значения string (ни int64, ни float64)",
			url:  "/update/counter/someMetric/some500Value",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run(test.name, func(t *testing.T) {
				request := httptest.NewRequest(http.MethodPost, test.url, nil)
				w := httptest.NewRecorder()

				CheckTypeAndValueMiddleware(getFakeHandler()).ServeHTTP(w, request)

				res := w.Result()
				defer res.Body.Close()

				assert.Equal(t, test.want.code, res.StatusCode)
			})
		})
	}
}
