package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/models"
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

func TestCheckGaugeType(t *testing.T) {
	testCases := []struct {
		name        string
		metricType  string
		metricValue string
		expected    bool
	}{
		{"корректный случай #1", models.GaugeType, "3.14", true},
		{"корректный случай #2", models.GaugeType, "42", true},
		{"некорректное значение", models.GaugeType, "foo", false},
		{"некорректный тип", models.CounterType, "42", false},
		{"пустые строки", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CheckGaugeType(tc.metricType, tc.metricValue)

			if actual != tc.expected {
				t.Errorf("CheckGaugeType(%q, %q) = %v; want %v", tc.metricType, tc.metricValue, actual, tc.expected)
			}
		})
	}
}

func TestCheckCounterType(t *testing.T) {
	testCases := []struct {
		name        string
		metricType  string
		metricValue string
		expected    bool
	}{
		{"корректный случай #1", models.CounterType, "3", true},
		{"корректный случай #2", models.CounterType, "42", true},
		{"некорректное значение", models.CounterType, "foo", false},
		{"некорректный тип", models.GaugeType, "42", false},
		{"пустые строки", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CheckCounterType(tc.metricType, tc.metricValue)

			if actual != tc.expected {
				t.Errorf("CheckCounterType(%q, %q) = %v; want %v", tc.metricType, tc.metricValue, actual, tc.expected)
			}
		})
	}
}

func TestSliceStrings(t *testing.T) {
	testCases := []struct {
		name     string
		strings  []string
		expected []string
		i        int
	}{
		{name: "корректный случай #1", strings: []string{"a", "b", "c"}, i: 1, expected: []string{"a", "c"}},
		{name: "корректный случай #2", strings: []string{"a", "b", "c"}, i: 0, expected: []string{"b", "c"}},
		{name: "корректный случай #3", strings: []string{"a", "b", "c"}, i: 2, expected: []string{"a", "b"}},
		{name: "индекс за пределами слайса", strings: []string{"a", "b", "c"}, i: 3, expected: []string{"a", "b", "c"}},
		{name: "отрицательный индекс", strings: []string{"a", "b", "c"}, i: -1, expected: []string{"a", "b", "c"}},
		{name: "пустой слайс", strings: []string{}, i: 0, expected: []string{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SliceStrings(tc.strings, tc.i)
			if !assert.ElementsMatch(t, actual, tc.expected) {
				t.Errorf("SliceStrings(%v, %d) = %v; want %v", tc.strings, tc.i, actual, tc.expected)
			}
		})
	}
}
