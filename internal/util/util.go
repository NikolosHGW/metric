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

const (
	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	Frees         = "Frees"
	GCCPUFraction = "GCCPUFraction"
	GCSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MCacheSys     = "MCacheSys"
	MSpanInuse    = "MSpanInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	TotalAlloc    = "TotalAlloc"
	PollCount     = "PollCount"
	RandomValue   = "RandomValue"
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

func GetMetricTypeMap() map[string]string {
	return map[string]string{
		Alloc:         GaugeType,
		BuckHashSys:   GaugeType,
		Frees:         GaugeType,
		GCCPUFraction: GaugeType,
		GCSys:         GaugeType,
		HeapAlloc:     GaugeType,
		HeapIdle:      GaugeType,
		HeapInuse:     GaugeType,
		HeapObjects:   GaugeType,
		HeapReleased:  GaugeType,
		HeapSys:       GaugeType,
		LastGC:        GaugeType,
		Lookups:       GaugeType,
		MCacheInuse:   GaugeType,
		MCacheSys:     GaugeType,
		MSpanInuse:    GaugeType,
		MSpanSys:      GaugeType,
		Mallocs:       GaugeType,
		NextGC:        GaugeType,
		NumForcedGC:   GaugeType,
		NumGC:         GaugeType,
		OtherSys:      GaugeType,
		PauseTotalNs:  GaugeType,
		StackInuse:    GaugeType,
		StackSys:      GaugeType,
		Sys:           GaugeType,
		TotalAlloc:    GaugeType,
		PollCount:     CounterType,
		RandomValue:   GaugeType,
	}
}
