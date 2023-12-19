package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/server/logger"

	"go.uber.org/zap"
)

type metricService interface {
	SetMetric(models.Metrics)
	GetMetricByName(string) (models.Metrics, error)
	GetAllMetrics() []string
}

type customLogger interface {
	Debug(string, ...zap.Field)
}

type Handler struct {
	metricService metricService
	logger        customLogger
}

func NewHandler(ms metricService, l customLogger) *Handler {
	return &Handler{
		metricService: ms,
		logger:        l,
	}
}

func (h Handler) SetMetric(w http.ResponseWriter, r *http.Request) {
	metricModel := models.NewMetricsModel()
	err := metricModel.DecodeMetricRequest(r.Body)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_SetMetric cannot decode request JSON body", zap.Error(err))
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.metricService.SetMetric(*metricModel)

	updatedMetric, err := h.metricService.GetMetricByName(metricModel.ID)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_SetMetric", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedMetric)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_SetMetric cannot encode to JSON", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricModel := models.NewMetricsModel()
	err := metricModel.DecodeMetricRequest(r.Body)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_GetMetric cannot decode request JSON body", zap.Error(err))
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	metric, err := h.metricService.GetMetricByName(metricModel.ID)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_GetMetric metric not found", zap.Error(err))
		http.Error(w, "метрика не найдена", http.StatusNotFound)

		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Debug("metric/internal/server/handlers/handlers.go Handler_GetMetric cannot encode to JSON", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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
