package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NikolosHGW/metric/internal/models"
)

type Repository interface {
	SetMetric(context.Context, models.Metrics) error
	GetMetric(context.Context, string) (models.Metrics, error)
	SetGaugeMetric(context.Context, string, models.Gauge) error
	SetCounterMetric(context.Context, string, models.Counter) error
	GetGaugeMetric(context.Context, string) (models.Gauge, error)
	GetCounterMetric(context.Context, string) (models.Counter, error)
	GetAllMetrics(context.Context) []string
	GetIsDBConnected() bool
	UpsertMetrics(context.Context, models.MetricCollection) (models.MetricCollection, error)
}

type MetricService struct {
	strg Repository
}

func NewMetricService(repo Repository) *MetricService {
	return &MetricService{
		strg: repo,
	}
}

func (ms MetricService) SetMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	var err error
	if metricType == models.CounterType {
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		err = ms.strg.SetCounterMetric(ctx, metricName, models.Counter(value))
	}

	if metricType == models.GaugeType {
		value, _ := strconv.ParseFloat(metricValue, 64)
		err = ms.strg.SetGaugeMetric(ctx, metricName, models.Gauge(value))
	}

	return err
}

func (ms MetricService) GetMetricValue(ctx context.Context, metricType, metricName string) (string, error) {
	if metricType == models.GaugeType {
		metricValue, err := ms.strg.GetGaugeMetric(ctx, metricName)

		return fmt.Sprintf("%v", metricValue), err
	}

	metricValue, err := ms.strg.GetCounterMetric(ctx, metricName)

	return fmt.Sprintf("%v", metricValue), err
}

func (ms *MetricService) SetJSONMetric(ctx context.Context, m models.Metrics) error {
	return ms.strg.SetMetric(ctx, m)
}

func (ms MetricService) GetMetricByName(ctx context.Context, name string) (models.Metrics, error) {
	return ms.strg.GetMetric(ctx, name)
}

func (ms MetricService) GetAllMetrics(ctx context.Context) []string {
	return ms.strg.GetAllMetrics(ctx)
}

func (ms MetricService) GetIsDBConnected() bool {
	return ms.strg.GetIsDBConnected()
}

func (ms MetricService) UpsertMetrics(ctx context.Context, mc models.MetricCollection) (models.MetricCollection, error) {
	return ms.strg.UpsertMetrics(ctx, mc)
}
