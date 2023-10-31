package storage

import "errors"

const (
	GAUGE_TYPE   = "gauge"
	COUNTER_TYPE = "counter"
)

type Metric struct {
	Gauge   float64
	Counter int64
}

type MemStorage struct {
	Metrics map[string]Metric
}

func (ms MemStorage) GetMetric(name string) (Metric, error) {
	metric, exist := ms.Metrics[name]
	if exist {
		return metric, nil
	}

	return Metric{}, errors.New("not found")
}

type MetricStorage interface {
	GetMetric(string)
}
