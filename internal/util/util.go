package util

import (
	"net/http"
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

type middleware func(http.Handler) http.Handler

func MiddlewareConveyor(handler http.Handler, middlewares ...middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

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
