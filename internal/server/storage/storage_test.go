package storage

import (
	"errors"
	"testing"

	"github.com/NikolosHGW/metric/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGetGaugeMetric(t *testing.T) {
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
