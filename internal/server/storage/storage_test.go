package storage

import (
	"errors"
	"testing"

	"github.com/NikolosHGW/metric/internal/util"
	"github.com/stretchr/testify/assert"
)

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
		{"отрицательный тест: достать несуществующую метрику baz", "baz", 0, errors.New("gauge metric not found")},
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
		{"отрицательный тест: достать несуществующую метрику baz", "baz", 0, errors.New("counter metric not found")},
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
