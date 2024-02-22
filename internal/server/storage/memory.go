package storage

import (
	"fmt"
	"sort"

	"github.com/NikolosHGW/metric/internal/models"
)

type metricValue struct {
	gauge   models.Gauge
	counter models.Counter
}

type MemStorage struct {
	metrics map[string]metricValue
}

func (ms MemStorage) GetGaugeMetric(name string) (models.Gauge, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.gauge, nil
	}

	return 0, fmt.Errorf("gauge metric %s not found", name)
}

func (ms MemStorage) GetCounterMetric(name string) (models.Counter, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.counter, nil
	}

	return 0, fmt.Errorf("counter metric %s not found", name)
}

func (ms *MemStorage) SetGaugeMetric(name string, value models.Gauge) {
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

func (ms *MemStorage) SetCounterMetric(name string, value models.Counter) {
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

func (ms *MemStorage) SetMetric(m models.Metrics) {
	if m.MType == models.CounterType {
		ms.SetCounterMetric(m.ID, models.Counter(*m.Delta))

		return
	}

	ms.SetGaugeMetric(m.ID, models.Gauge(*m.Value))
}

func getMetricsModel(name string, metric metricValue) models.Metrics {
	if metric.counter != 0 {
		return models.Metrics{ID: name, MType: models.CounterType, Delta: (*int64)(&metric.counter)}
	}

	return models.Metrics{ID: name, MType: models.GaugeType, Value: (*float64)(&metric.gauge)}
}

func (ms *MemStorage) GetMetric(name string) (models.Metrics, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return getMetricsModel(name, metric), nil
	}

	return models.Metrics{}, fmt.Errorf("%s metric not found", name)
}

func (ms *MemStorage) GetMetricsModels() []models.Metrics {
	models := make([]models.Metrics, 0, len(ms.metrics))
	for k := range ms.metrics {
		model, _ := ms.GetMetric(k)
		models = append(models, model)
	}

	return models
}

func (ms MemStorage) GetAllMetrics() []string {
	result := make([]string, len(ms.metrics))

	keys := make([]string, 0, len(ms.metrics))
	for k := range ms.metrics {
		keys = append(keys, k)
	}

	// Сортируем срез ключей
	sort.Strings(keys)

	// Итерируем по отсортированному срезу и заполняем результат
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
