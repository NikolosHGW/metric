package memory

import (
	"fmt"
	"sort"

	"github.com/NikolosHGW/metric/internal/util"
)

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

	return 0, fmt.Errorf("gauge metric %s not found", name)
}

func (ms MemStorage) GetCounterMetric(name string) (util.Counter, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric.counter, nil
	}

	return 0, fmt.Errorf("counter metric %s not found", name)
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
	return new(MemStorage)
}
