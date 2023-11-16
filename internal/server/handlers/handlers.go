package handlers

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/services/metric"
	"github.com/NikolosHGW/metric/internal/server/storage"
)

func WithSetMetricHandle(strg storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric.SetMetric(r, strg)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func WithGetValueMetricHandle(strg storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricValue := metric.GetMetricValue()

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(metricValue))
	}
}
