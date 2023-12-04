package metric

import (
	"fmt"
	"strconv"

	"github.com/NikolosHGW/metric/internal/util"
)

type repository interface {
	SetGaugeMetric(string, util.Gauge)
	SetCounterMetric(string, util.Counter)
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

func (ms MetricService) SetMetric(metricType, metricName, metricValue string) {
	if metricType == util.CounterType {
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		ms.strg.SetCounterMetric(metricName, util.Counter(value))
	}

	if metricType == util.GaugeType {
		value, _ := strconv.ParseFloat(metricValue, 64)
		ms.strg.SetGaugeMetric(metricName, util.Gauge(value))
	}
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
