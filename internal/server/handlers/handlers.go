package handlers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
)

type metricService interface {
	SetMetric(string, string, string)
	GetMetricValue(string, string) (string, error)
	GetAllMetrics() []string
}

type Handler struct {
	metricService metricService
}

func NewHandler(ms metricService) *Handler {
	return &Handler{
		metricService: ms,
	}
}

func (h Handler) SetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	h.metricService.SetMetric(metricType, metricName, metricValue)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Content-Type", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h Handler) GetValueMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue, err := h.metricService.GetMetricValue(metricType, metricName)
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
	metrics := h.metricService.GetAllMetrics()

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
