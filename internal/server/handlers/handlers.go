package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type metricService interface {
	SetMetric(models.Metrics)
	GetMetricByName(string) (models.Metrics, error)
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
	metricModel := models.NewMetricsModel()
	err := metricModel.DecodeMetricRequest(r.Body)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_SetMetric cannot decode request JSON body", zap.Error(err))
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.metricService.SetMetric(*metricModel)

	updatedMetric, err := h.metricService.GetMetricByName(metricModel.ID)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_SetMetric", zap.Error(err))
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedMetric)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_SetMetric cannot encode to JSON", zap.Error(err))
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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
