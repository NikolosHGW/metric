package memory

import (
	"testing"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetMetric(t *testing.T) {
	ms := NewMemStorage()
	fooValue := 42.1
	mockModel := models.Metrics{
		ID:    "foo",
		Value: &fooValue,
	}

	ms.SetMetric(mockModel)

	metric, exist := ms.metrics["foo"]
	assert.True(t, exist, "метрика не найдена в хранилище")
	assert.Equal(t, mockModel, metric, "метрика не соответствует установленной")
}

func TestMemStorage_GetMetric(t *testing.T) {
	fooValue := 42.1
	var barValue int64 = 100
	ms := &MemStorage{
		metrics: map[string]models.Metrics{
			"foo": {
				ID:    "foo",
				MType: util.GaugeType,
				Value: &fooValue,
			},
			"bar": {
				ID:    "bar",
				MType: util.CounterType,
				Delta: &barValue,
			},
		},
	}

	testCases := []struct {
		name     string
		metric   string
		expected models.Metrics
		err      bool
	}{
		{
			name:     "достать существующую gauge метрику",
			metric:   "foo",
			expected: ms.metrics["foo"],
			err:      false,
		},
		{
			name:     "достать существующую counter метрику",
			metric:   "bar",
			expected: ms.metrics["bar"],
			err:      false,
		},
		{
			name:     "достать несуществующую метрику",
			metric:   "baz",
			expected: models.Metrics{},
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actual, err := ms.GetMetric(tc.metric)

			assert.Equal(t, tc.expected, actual)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	fooValue := 42.1
	mockModel1 := models.Metrics{
		ID:    "foo",
		Value: &fooValue,
	}
	barValue := 100.01
	mockModel2 := models.Metrics{
		ID:    "bar",
		Value: &barValue,
	}

	testCases := []struct {
		name     string
		input    []models.Metrics
		expected []string
	}{
		{
			name:     "положительный тест: с наполненным сторэджом",
			input:    []models.Metrics{mockModel1, mockModel2},
			expected: []string{"bar: 100.01", "foo: 42.1"},
		},
		{
			name:     "отрицательный тест: пустой сторэдж",
			input:    []models.Metrics{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := NewMemStorage()
			for _, model := range tc.input {
				ms.SetMetric(model)
			}
			actual := ms.GetAllMetrics()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
