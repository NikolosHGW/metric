package metric

import (
	"fmt"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/util"
)

type repository interface {
	SetMetric(models.Metrics)
	GetMetric(string) (models.Metrics, error)
	GetGaugeMetric(string) (util.Gauge, error)
	GetCounterMetric(string) (util.Counter, error)
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

func (ms MetricService) SetMetric(m models.Metrics) {
	ms.strg.SetMetric(m)
}

func (ms MetricService) GetMetricByName(name string) (models.Metrics, error) {
	return ms.strg.GetMetric(name)
}

func (ms MetricService) GetMetricValue(metricType, metricName string) (string, error) {
	if metricType == util.GaugeType {
		metricValue, err := ms.strg.GetGaugeMetric(metricName)

		return fmt.Sprintf("%v", metricValue), err
	}

	metricValue, err := ms.strg.GetCounterMetric(metricName)

	return fmt.Sprintf("%v", metricValue), err
}

func (ms MetricService) GetAllMetrics() []string {
	return ms.strg.GetAllMetrics()
}
