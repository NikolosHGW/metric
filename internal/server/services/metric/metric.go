package metric

import (
	"fmt"
	"strconv"

	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/NikolosHGW/metric/internal/util"
)

func SetMetric(strg storage.Storage, metricType, metricName, metricValue string) {
	if metricType == util.CounterType {
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		strg.SetCounterMetric(metricName, util.Counter(value))
	}

	if metricType == util.GaugeType {
		value, _ := strconv.ParseFloat(metricValue, 64)
		strg.SetGaugeMetric(metricName, util.Gauge(value))
	}
}

func GetMetricValue(strg storage.Storage, metricType, metricName string) (string, error) {
	if metricType == util.GaugeType {
		metricValue, err := strg.GetGaugeMetric(metricName)

		return fmt.Sprintf("%v", metricValue), err
	}

	metricValue, err := strg.GetCounterMetric(metricName)

	return fmt.Sprintf("%v", metricValue), err
}
