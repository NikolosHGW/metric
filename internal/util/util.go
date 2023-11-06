package util

import (
	"strconv"
)

const (
	Update = iota
	MetricType
	MetricName
	MetricValue
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Gauge float64
type Counter int64

func CheckGaugeType(metricType string, metricValue string) bool {
	_, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return false
	}

	return metricType == GaugeType
}

func CheckCounterType(metricType string, metricValue string) bool {
	_, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return false
	}

	return metricType == CounterType
}

func SliceStrings(strings []string, i int) []string {
	if len(strings) != 0 && i < len(strings) {
		strings = append(strings[:i], strings[i+1:]...)
	}

	return strings
}
