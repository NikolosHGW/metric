package handlers

import (
	"html/template"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/services/metric"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/go-chi/chi"
)

func WithSetMetricHandle(strg storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")
		metric.SetMetric(strg, metricType, metricName, metricValue)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func WithGetValueMetricHandle(strg storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue, err := metric.GetMetricValue(strg, metricType, metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(metricValue))
	}
}

func WithGetMetricsHandle(strg storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := strg.GetAllMetrics()

		tmpl, err := template.ParseFiles("../../internal/server/templates/list_metrics.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
