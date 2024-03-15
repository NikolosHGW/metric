package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/go-chi/chi"

	"go.uber.org/zap"
)

type metricService interface {
	SetMetric(string, string, string, context.Context) error
	SetJSONMetric(models.Metrics, context.Context) error
	GetMetricValue(string, string, context.Context) (string, error)
	GetMetricByName(string, context.Context) (models.Metrics, error)
	GetAllMetrics(context.Context) []string
	GetIsDBConnected() bool
	UpsertMetrics(models.MetricCollection, context.Context) (models.MetricCollection, error)
}

type customLogger interface {
	Info(string, ...zap.Field)
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
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	err := h.metricService.SetMetric(metricType, metricName, metricValue, r.Context())
	if err != nil {
		h.logger.Info("cannot upsert metric", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Content-Type", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h Handler) GetValueMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue, err := h.metricService.GetMetricValue(metricType, metricName, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Content-Type", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricValue))
}

func (h Handler) SetJSONMetric(w http.ResponseWriter, r *http.Request) {
	metricModel := models.NewMetricModel()
	err := metricModel.DecodeMetricRequest(r.Body)
	if err != nil {
		h.logger.Info("cannot decode request JSON body", zap.Error(err))
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = h.metricService.SetJSONMetric(*metricModel, r.Context())
	if err != nil {
		h.logger.Info("cannot upsert metric", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	updatedMetric, err := h.metricService.GetMetricByName(metricModel.ID, r.Context())
	if err != nil {
		h.logger.Info("ошибка в GetMetricByName", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedMetric)
	if err != nil {
		h.logger.Info("cannot encode to JSON", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Println("in: ", r.Body, "; ", "out: ", resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricModel := models.NewMetricModel()
	err := metricModel.DecodeMetricRequest(r.Body)
	if err != nil {
		h.logger.Info("cannot decode request JSON body", zap.Error(err))
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	metric, err := h.metricService.GetMetricByName(metricModel.ID, r.Context())
	if err != nil {
		h.logger.Info("metric not found", zap.Error(err))
		http.Error(w, "метрика не найдена", http.StatusNotFound)

		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		h.logger.Info("cannot encode to JSON", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Println("in: ", r.Body, "; ", "out: ", resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.metricService.GetAllMetrics(r.Context())

	tmpl, err := template.ParseFiles(filepath.Join(BasePath(), "/internal/server/handlers/list_metrics.tmpl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func BasePath() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "../../..")
}

func (h Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	if h.metricService.GetIsDBConnected() {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (h Handler) UpsertMetrics(w http.ResponseWriter, r *http.Request) {
	metricCollection := models.NewMetricCollection()
	err := metricCollection.DecodeMetricsRequest(r.Body)
	if err != nil {
		h.logger.Info("cannot decode request JSON body", zap.Error(err))
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	metrics, err := h.metricService.UpsertMetrics(*metricCollection, r.Context())
	if err != nil {
		h.logger.Info("cannot upsert metrics", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(metrics.Metrics)
	if err != nil {
		h.logger.Info("cannot encode to JSON", zap.Error(err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Println("in: ", r.Body, "; ", "out: ", resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
