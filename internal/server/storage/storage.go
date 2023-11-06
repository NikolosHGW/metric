package storage

import (
	"errors"

	"github.com/NikolosHGW/metric/internal/util"
)

type MetricStorage interface {
	GetGaugeMetric(string) (util.Gauge, error)
	GetCounterMetric(string) (util.Counter, error)
	SetGaugeMetric(string, util.Gauge)
	SetCounterMetric(string, util.Counter)
}

type metricValue struct {
	gauge   util.Gauge
	counter util.Counter
}

type MemStorage struct {
	metrics map[string]metricValue
}

func (ms MemStorage) GetGaugeMetric(name string) (util.Gauge, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.gauge, nil
	}

	return 0, errors.New("not found")
}

func (ms MemStorage) GetCounterMetric(name string) (util.Counter, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.counter, nil
	}

	return 0, errors.New("not found")
}

func (ms *MemStorage) SetGaugeMetric(name string, value util.Gauge) {
	metric, exist := ms.metrics[name]
	if exist {
		metric.gauge = value
		ms.metrics[name] = metric
	} else {
		if ms.metrics == nil {
			ms.metrics = make(map[string]metricValue)
		}
		ms.metrics[name] = metricValue{
			gauge: value,
		}
	}
}

func (ms *MemStorage) SetCounterMetric(name string, value util.Counter) {
	metric, exist := ms.metrics[name]
	if exist {
		metric.counter += value
		ms.metrics[name] = metric
	} else {
		if ms.metrics == nil {
			ms.metrics = make(map[string]metricValue)
		}
		ms.metrics[name] = metricValue{
			counter: value,
		}
	}
}

func NewMemStorage() *MemStorage {
	return new(MemStorage)
}
