package metric

import (
	"github.com/NikolosHGW/metric/internal/models"
)

type repository interface {
	SetMetric(models.Metrics)
	GetMetric(string) (models.Metrics, error)
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

func (ms MetricService) GetAllMetrics() []string {
	return ms.strg.GetAllMetrics()
}
