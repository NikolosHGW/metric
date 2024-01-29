package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/NikolosHGW/metric/internal/models"
)

const (
	Update = iota
	MetricType
	MetricName
	MetricValue
)

func SliceStrings(strings []string, i int) []string {
	if len(strings) != 0 && i < len(strings) && i >= 0 {
		strings = append(strings[:i], strings[i+1:]...)
	}

	return strings
}

func CheckGaugeType(metricType string, metricValue string) bool {
	_, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return false
	}

	return metricType == models.GaugeType
}

func CheckCounterType(metricType string, metricValue string) bool {
	_, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return false
	}

	return metricType == models.CounterType
}

func CheckMetricNameMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := SliceStrings(strings.Split(r.URL.Path, "/"), 0)

		if len(parts) != 4 {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckTypeAndValueMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := SliceStrings(strings.Split(r.URL.Path, "/"), 0)

		isCounterType := CheckCounterType(parts[MetricType], parts[MetricValue])
		isGaugeType := CheckGaugeType(parts[MetricType], parts[MetricValue])

		if !isCounterType && !isGaugeType {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}
