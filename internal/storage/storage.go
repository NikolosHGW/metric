package storage

import (
	"errors"
)

type Gauge float64
type Counter int64

type MetricStorage interface {
	GetGaugeMetric(string) (Gauge, error)
	GetCounterMetric(string) (Counter, error)
	SetGaugeMetric(string, Gauge)
	SetCounterMetric(string, Counter)
}

type metricValue struct {
	gauge   Gauge
	counter Counter
}

type MemStorage struct {
	metrics map[string]metricValue
}

func (ms MemStorage) GetGaugeMetric(name string) (Gauge, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.gauge, nil
	}

	return 0, errors.New("not found")
}

func (ms MemStorage) GetCounterMetric(name string) (Counter, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.counter, nil
	}

	return 0, errors.New("not found")
}

func (ms *MemStorage) SetGaugeMetric(name string, value Gauge) {
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

func (ms *MemStorage) SetCounterMetric(name string, value Counter) {
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

func InitStorage() MetricStorage {
	return &MemStorage{}
}
