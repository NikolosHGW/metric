package memory

import (
	"errors"
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
	assert.Equal(t, util.Gauge(*mockModel.Value), metric.gauge, "метрика не соответствует установленной")
}

func TestMemStorage_GetMetric(t *testing.T) {
	fooValue := 42.1
	var barValue int64 = 100
	expectedMetric := map[string]models.Metrics{
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
	}
	ms := &MemStorage{
		metrics: map[string]metricValue{
			"foo": {gauge: util.Gauge(fooValue)},
			"bar": {counter: util.Counter(barValue)},
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
			expected: expectedMetric["foo"],
			err:      false,
		},
		{
			name:     "достать существующую counter метрику",
			metric:   "bar",
			expected: expectedMetric["bar"],
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

func TestMemStorage_GetGaugeMetric(t *testing.T) {
	ms := NewMemStorage()
	ms.SetGaugeMetric("foo", 42.1)
	ms.SetGaugeMetric("bar", 100.01)

	testCases := []struct {
		name       string
		metricName string
		expected   util.Gauge
		err        error
	}{
		{"положительный тест: достать существующую метрику foo", "foo", 42.1, nil},
		{"положительный тест: достать существующую метрику bar", "bar", 100.01, nil},
		{"отрицательный тест: достать несуществующую метрику baz", "baz", 0, errors.New("gauge metric baz not found")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ms.GetGaugeMetric(tc.metricName)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestMemStorage_GetCounterMetric(t *testing.T) {
	ms := NewMemStorage()
	ms.SetCounterMetric("foo", 42)
	ms.SetCounterMetric("bar", 100)

	testCases := []struct {
		name       string
		metricName string
		expected   util.Counter
		err        error
	}{
		{"положительный тест: достать существующую метрику foo", "foo", 42, nil},
		{"положительный тест: достать существующую метрику bar", "bar", 100, nil},
		{"отрицательный тест: достать несуществующую метрику baz", "baz", 0, errors.New("counter metric baz not found")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ms.GetCounterMetric(tc.metricName)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestMemStorage_SetGaugeMetric(t *testing.T) {
	ms := NewMemStorage()

	testCases := []struct {
		name       string
		metricName string
		value      util.Gauge
	}{
		{"положительный тест: устаналиваем метрику foo", "foo", 42.1},
		{"положительный тест: устаналиваем метрику bar", "bar", 100.01},
		{"положительный тест: устанавливаем метрику пустую строку и отрицательное значение", "", -20.2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms.SetGaugeMetric(tc.metricName, tc.value)
			actual, err := ms.GetGaugeMetric(tc.metricName)
			assert.Equal(t, tc.value, actual)
			assert.Nil(t, err)
		})
	}
}

func TestMemStorage_SetCounterMetric(t *testing.T) {
	ms := NewMemStorage()

	testCases := []struct {
		name       string
		metricName string
		value      util.Counter
	}{
		{"положительный тест: устаналиваем метрику foo", "foo", 42},
		{"положительный тест: устаналиваем метрику bar", "bar", 100},
		{"положительный тест: устанавливаем метрику пустую строку и отрицательное значение", "", -20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms.SetCounterMetric(tc.metricName, tc.value)
			actual, err := ms.GetCounterMetric(tc.metricName)
			assert.Equal(t, tc.value, actual)
			assert.Nil(t, err)
		})
	}
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	type data struct {
		metricName  string
		metricValue util.Gauge
	}

	testCases := []struct {
		name     string
		input    []data
		expected []string
	}{
		{
			name:     "положительный тест: с наполненным сторэджом",
			input:    []data{{metricName: "foo", metricValue: 42}, {metricName: "bar", metricValue: 100}},
			expected: []string{"bar: 100", "foo: 42"},
		},
		{
			name:     "отрицательный тест: пустой сторэдж",
			input:    []data{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := NewMemStorage()
			for _, data := range tc.input {
				ms.SetGaugeMetric(data.metricName, data.metricValue)
			}
			actual := ms.GetAllMetrics()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
