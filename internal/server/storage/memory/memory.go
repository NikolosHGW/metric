package memory

import (
	"fmt"
	"sort"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/util"
)

type MemStorage struct {
	metrics map[string]models.Metrics
}

func (ms *MemStorage) SetMetric(m models.Metrics) {
	ms.metrics[m.ID] = m
}

func (ms *MemStorage) GetMetric(name string) (models.Metrics, error) {
	metric, exist := ms.metrics[name]
	if exist {
		return metric, nil
	}

	return models.Metrics{}, fmt.Errorf("%s metric not found", name)
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
		if v.MType == util.CounterType {
			result[i] = fmt.Sprintf("%v: %v", k, *v.Delta)
		} else {
			result[i] = fmt.Sprintf("%v: %v", k, *v.Value)
		}
		i++
	}

	return result
}

func NewMemStorage() *MemStorage {
	storage := new(MemStorage)
	storage.metrics = make(map[string]models.Metrics)

	return storage
}
