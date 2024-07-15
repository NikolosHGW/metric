package storage

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/NikolosHGW/metric/internal/models"
)

type metricValue struct {
	gauge   models.Gauge
	counter models.Counter
}

type MemStorage struct {
	metrics map[string]metricValue
	mtx     sync.Mutex
}

func (ms *MemStorage) GetGaugeMetric(_ context.Context, name string) (models.Gauge, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.gauge, nil
	}

	return 0, fmt.Errorf("gauge metric %s not found", name)
}

func (ms *MemStorage) GetCounterMetric(_ context.Context, name string) (models.Counter, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.counter, nil
	}

	return 0, fmt.Errorf("counter metric %s not found", name)
}

func (ms *MemStorage) SetGaugeMetric(_ context.Context, name string, value models.Gauge) error {
	ms.mtx.Lock()
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

	ms.mtx.Unlock()

	return nil
}

func (ms *MemStorage) SetCounterMetric(_ context.Context, name string, value models.Counter) error {
	ms.mtx.Lock()
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

	ms.mtx.Unlock()

	return nil
}

func (ms *MemStorage) SetMetric(ctx context.Context, m models.Metrics) error {
	if m.MType == models.CounterType {
		err := ms.SetCounterMetric(ctx, m.ID, models.Counter(*m.Delta))
		if err != nil {
			return fmt.Errorf("can not SetCounterMetric: %w", err)
		}

		return nil
	}

	err := ms.SetGaugeMetric(ctx, m.ID, models.Gauge(*m.Value))
	if err != nil {
		return fmt.Errorf("can not SetGaugeMetric: %w", err)
	}

	return nil
}

func getMetricsModel(_ context.Context, name string, metric metricValue) models.Metrics {
	if metric.counter != 0 {
		return models.Metrics{ID: name, MType: models.CounterType, Delta: (*int64)(&metric.counter)}
	}

	return models.Metrics{ID: name, MType: models.GaugeType, Value: (*float64)(&metric.gauge)}
}

func (ms *MemStorage) GetMetric(ctx context.Context, name string) (models.Metrics, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return getMetricsModel(ctx, name, metric), nil
	}

	return models.Metrics{}, fmt.Errorf("%s metric not found", name)
}

func (ms *MemStorage) GetMetricsModels(ctx context.Context) []models.Metrics {
	models := make([]models.Metrics, 0, len(ms.metrics))
	for k := range ms.metrics {
		model, _ := ms.GetMetric(ctx, k)
		models = append(models, model)
	}

	return models
}

func (ms *MemStorage) GetAllMetrics(_ context.Context) []string {
	result := make([]string, len(ms.metrics))

	keys := make([]string, 0, len(ms.metrics))
	for k := range ms.metrics {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	i := 0
	for _, k := range keys {
		v := ms.metrics[k]
		if v.counter != 0 {
			result[i] = fmt.Sprintf("%v: %v", k, v.counter)
		} else {
			result[i] = fmt.Sprintf("%v: %v", k, v.gauge)
		}
		i++
	}

	return result
}

func NewMemStorage() *MemStorage {
	storage := new(MemStorage)
	storage.metrics = make(map[string]metricValue)

	return storage
}

func (ms *MemStorage) GetIsDBConnected() bool {
	return false
}

func (ms *MemStorage) UpsertMetrics(ctx context.Context, metricCollection models.MetricCollection) (models.MetricCollection, error) {
	for _, m := range metricCollection.Metrics {
		err := ms.SetMetric(ctx, m)
		if err != nil {
			return metricCollection, fmt.Errorf("can not SetMetric: %w", err)
		}
	}

	return metricCollection, nil
}
