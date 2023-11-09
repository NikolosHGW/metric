package util

import (
	"testing"

	"github.com/NikolosHGW/metric/internal/util"
)

func Test_getStringValue(t *testing.T) {
	testCases := []struct {
		name     string
		v        interface{}
		expected string
	}{
		{"корректный случай #1", util.Gauge(3.14), "3.14"},
		{"корректный случай #2", util.Counter(42), "42"},
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
		metricType  string
		metricName  string
		metricValue string
		expected    string
	}{
		{"корректный случай #1", util.GaugeType, "Alloc", "123.45", "http://localhost:8080/update/gauge/Alloc/123.45"},
		{"корректный случай #2", util.CounterType, "PollCount", "42", "http://localhost:8080/update/counter/PollCount/42"},
		{"пустой тип метрики", "", "Alloc", "123.45", "http://localhost:8080/update/Alloc/123.45"},
		{"пустое имя метрики", util.GaugeType, "", "123.45", "http://localhost:8080/update/gauge/123.45"},
		{"пустое значение метрики", util.CounterType, "PollCounter", "", "http://localhost:8080/update/counter/PollCounter/"},
		{"пустые строки", "", "", "", "http://localhost:8080/update/"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getResultUrl(tc.metricType, tc.metricName, tc.metricValue)
			if actual != tc.expected {
				t.Errorf("getResultUrl(%q, %q, %q) = %q; want %q", tc.metricType, tc.metricName, tc.metricValue, actual, tc.expected)
			}
		})
	}
}
