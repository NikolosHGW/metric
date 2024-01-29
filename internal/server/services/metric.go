package services

import (
	"fmt"
	"strconv"

	"github.com/NikolosHGW/metric/internal/models"
)

type repository interface {
	SetMetric(models.Metrics)
	GetMetric(string) (models.Metrics, error)
	SetGaugeMetric(string, models.Gauge)
	SetCounterMetric(string, models.Counter)
	GetGaugeMetric(string) (models.Gauge, error)
	GetCounterMetric(string) (models.Counter, error)
	GetAllMetrics() []string
}

type MetricService struct {
	strg repository
}

func NewMetricService(repo repository) *MetricService {
	return &MetricService{
		strg: repo,
	}
}

func (ms MetricService) SetMetric(metricType, metricName, metricValue string) {
	if metricType == models.CounterType {
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		ms.strg.SetCounterMetric(metricName, models.Counter(value))
	}

	if metricType == models.GaugeType {
		value, _ := strconv.ParseFloat(metricValue, 64)
		ms.strg.SetGaugeMetric(metricName, models.Gauge(value))
	}
}

func (ms MetricService) GetMetricValue(metricType, metricName string) (string, error) {
	if metricType == models.GaugeType {
		metricValue, err := ms.strg.GetGaugeMetric(metricName)

		return fmt.Sprintf("%v", metricValue), err
	}

	metricValue, err := ms.strg.GetCounterMetric(metricName)

	return fmt.Sprintf("%v", metricValue), err
}

func (ms *MetricService) SetJSONMetric(m models.Metrics) {
	ms.strg.SetMetric(m)
}

func (ms MetricService) GetMetricByName(name string) (models.Metrics, error) {
	return ms.strg.GetMetric(name)
}

func (ms MetricService) GetAllMetrics() []string {
	return ms.strg.GetAllMetrics()
}
