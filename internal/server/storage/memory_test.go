package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetMetric(t *testing.T) {
	ms := NewMemStorage()
	fooValue := 42.1
	mockModel := models.Metrics{
		ID:    "foo",
		Value: &fooValue,
	}

	err := ms.SetMetric(context.Background(), mockModel)
	if err != nil {
		t.Errorf("failed to SetMetric: %v", err)
	}

	metric, exist := ms.metrics["foo"]
	assert.True(t, exist, "метрика не найдена в хранилище")
	assert.Equal(t, models.Gauge(*mockModel.Value), metric.gauge, "метрика не соответствует установленной")
}

func TestMemStorage_GetMetric(t *testing.T) {
	fooValue := 42.1
	var barValue int64 = 100
	expectedMetric := map[string]models.Metrics{
		"foo": {
			ID:    "foo",
			MType: models.GaugeType,
			Value: &fooValue,
		},
		"bar": {
			ID:    "bar",
			MType: models.CounterType,
			Delta: &barValue,
		},
	}
	ms := &MemStorage{
		metrics: map[string]metricValue{
			"foo": {gauge: models.Gauge(fooValue)},
			"bar": {counter: models.Counter(barValue)},
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

			actual, err := ms.GetMetric(context.Background(), tc.metric)

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
	err := ms.SetGaugeMetric(context.Background(), "foo", 42.1)
	if err != nil {
		t.Errorf("failed to SetGaugeMetric: %v", err)
	}
	err = ms.SetGaugeMetric(context.Background(), "bar", 100.01)
	if err != nil {
		t.Errorf("2 failed to SetGaugeMetric: %v", err)
	}

	testCases := []struct {
		err        error
		metricName string
		name       string
		expected   models.Gauge
	}{
		{name: "положительный тест: достать существующую метрику foo", metricName: "foo", expected: 42.1, err: nil},
		{name: "положительный тест: достать существующую метрику bar", metricName: "bar", expected: 100.01, err: nil},
		{name: "отрицательный тест: достать несуществующую метрику baz", metricName: "baz", expected: 0, err: errors.New("gauge metric baz not found")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ms.GetGaugeMetric(context.Background(), tc.metricName)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestMemStorage_GetCounterMetric(t *testing.T) {
	ms := NewMemStorage()
	err := ms.SetCounterMetric(context.Background(), "foo", 42)
	if err != nil {
		t.Errorf("failed to SetCounterMetric: %v", err)
	}
	err = ms.SetCounterMetric(context.Background(), "bar", 100)
	if err != nil {
		t.Errorf("2 failed to SetCounterMetric: %v", err)
	}

	testCases := []struct {
		err        error
		metricName string
		name       string
		expected   models.Counter
	}{
		{name: "положительный тест: достать существующую метрику foo", metricName: "foo", expected: 42, err: nil},
		{name: "положительный тест: достать существующую метрику bar", metricName: "bar", expected: 100, err: nil},
		{name: "отрицательный тест: достать несуществующую метрику baz", metricName: "baz", expected: 0, err: errors.New("counter metric baz not found")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ms.GetCounterMetric(context.Background(), tc.metricName)
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
		value      models.Gauge
	}{
		{"положительный тест: устаналиваем метрику foo", "foo", 42.1},
		{"положительный тест: устаналиваем метрику bar", "bar", 100.01},
		{"положительный тест: устанавливаем метрику пустую строку и отрицательное значение", "", -20.2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ms.SetGaugeMetric(context.Background(), tc.metricName, tc.value)
			if err != nil {
				t.Errorf("failed to SetGaugeMetric: %v", err)
			}
			actual, err := ms.GetGaugeMetric(context.Background(), tc.metricName)
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
		value      models.Counter
	}{
		{"положительный тест: устаналиваем метрику foo", "foo", 42},
		{"положительный тест: устаналиваем метрику bar", "bar", 100},
		{"положительный тест: устанавливаем метрику пустую строку и отрицательное значение", "", -20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ms.SetCounterMetric(context.Background(), tc.metricName, tc.value)
			if err != nil {
				t.Errorf("failed to SetCounterMetric: %v", err)
			}
			actual, err := ms.GetCounterMetric(context.Background(), tc.metricName)
			assert.Equal(t, tc.value, actual)
			assert.Nil(t, err)
		})
	}
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	type data struct {
		metricName  string
		metricValue models.Gauge
	}

	testCases := []struct {
		name     string
		expected []string
		input    []data
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
				err := ms.SetGaugeMetric(context.Background(), data.metricName, data.metricValue)
				if err != nil {
					t.Errorf("failed to SetGaugeMetric: %v", err)
				}
			}
			actual := ms.GetAllMetrics(context.Background())
			assert.Equal(t, tc.expected, actual)
		})
	}
}
