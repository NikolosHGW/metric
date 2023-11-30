package handlers

import (
	"html/template"
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/services/metric"
	"github.com/NikolosHGW/metric/internal/util"
	"github.com/go-chi/chi"
)

type Repository interface {
	SetGaugeMetric(string, util.Gauge)
	SetCounterMetric(string, util.Counter)
	GetGaugeMetric(string) (util.Gauge, error)
	GetCounterMetric(string) (util.Counter, error)
	GetAllMetrics() []string
}

type Handler struct {
	repo Repository
}

func NewHandler(r Repository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h Handler) SetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	metric.SetMetric(h.repo, metricType, metricName, metricValue)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Content-Type", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h Handler) GetValueMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue, err := metric.GetMetricValue(h.repo, metricType, metricName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Content-Type", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricValue))
}

func (h Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.repo.GetAllMetrics()

	tmpl, err := template.ParseFiles("list_metrics.tmpl")
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
