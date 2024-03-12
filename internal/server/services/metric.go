package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NikolosHGW/metric/internal/models"
)

type repository interface {
	SetMetric(models.Metrics, context.Context) error
	GetMetric(string, context.Context) (models.Metrics, error)
	SetGaugeMetric(string, models.Gauge, context.Context) error
	SetCounterMetric(string, models.Counter, context.Context) error
	GetGaugeMetric(string, context.Context) (models.Gauge, error)
	GetCounterMetric(string, context.Context) (models.Counter, error)
	GetAllMetrics(context.Context) []string
	GetIsDBConnected() bool
	UpsertMetrics(models.MetricCollection, context.Context) (models.MetricCollection, error)
}

type MetricService struct {
	strg repository
}

func NewMetricService(repo repository) *MetricService {
	return &MetricService{
		strg: repo,
	}
}

func (ms MetricService) SetMetric(metricType, metricName, metricValue string, ctx context.Context) error {
	var err error
	if metricType == models.CounterType {
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		err = ms.strg.SetCounterMetric(metricName, models.Counter(value), ctx)
	}

	if metricType == models.GaugeType {
		value, _ := strconv.ParseFloat(metricValue, 64)
		err = ms.strg.SetGaugeMetric(metricName, models.Gauge(value), ctx)
	}

	return err
}

func (ms MetricService) GetMetricValue(metricType, metricName string, ctx context.Context) (string, error) {
	if metricType == models.GaugeType {
		metricValue, err := ms.strg.GetGaugeMetric(metricName, ctx)

		return fmt.Sprintf("%v", metricValue), err
	}

	metricValue, err := ms.strg.GetCounterMetric(metricName, ctx)

	return fmt.Sprintf("%v", metricValue), err
}

func (ms *MetricService) SetJSONMetric(m models.Metrics, ctx context.Context) error {
	return ms.strg.SetMetric(m, ctx)
}

func (ms MetricService) GetMetricByName(name string, ctx context.Context) (models.Metrics, error) {
	return ms.strg.GetMetric(name, ctx)
}

func (ms MetricService) GetAllMetrics(ctx context.Context) []string {
	return ms.strg.GetAllMetrics(ctx)
}

func (ms MetricService) GetIsDBConnected() bool {
	return ms.strg.GetIsDBConnected()
}

func (ms MetricService) UpsertMetrics(mc models.MetricCollection, ctx context.Context) (models.MetricCollection, error) {
	return ms.strg.UpsertMetrics(mc, ctx)
}
