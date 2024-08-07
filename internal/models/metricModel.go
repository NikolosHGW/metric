package models

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Gauge float64
type Counter int64

const (
	Alloc           = "Alloc"
	BuckHashSys     = "BuckHashSys"
	Frees           = "Frees"
	GCCPUFraction   = "GCCPUFraction"
	GCSys           = "GCSys"
	HeapAlloc       = "HeapAlloc"
	HeapIdle        = "HeapIdle"
	HeapInuse       = "HeapInuse"
	HeapObjects     = "HeapObjects"
	HeapReleased    = "HeapReleased"
	HeapSys         = "HeapSys"
	LastGC          = "LastGC"
	Lookups         = "Lookups"
	MCacheInuse     = "MCacheInuse"
	MCacheSys       = "MCacheSys"
	MSpanInuse      = "MSpanInuse"
	MSpanSys        = "MSpanSys"
	Mallocs         = "Mallocs"
	NextGC          = "NextGC"
	NumForcedGC     = "NumForcedGC"
	NumGC           = "NumGC"
	OtherSys        = "OtherSys"
	PauseTotalNs    = "PauseTotalNs"
	StackInuse      = "StackInuse"
	StackSys        = "StackSys"
	Sys             = "Sys"
	TotalAlloc      = "TotalAlloc"
	PollCount       = "PollCount"
	RandomValue     = "RandomValue"
	TotalMemory     = "TotalMemory"
	FreeMemory      = "FreeMemory"
	CPUutilization1 = "CPUutilization1"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	ID    string   `json:"id" db:"id"`                 // имя метрики
	MType string   `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
}

func NewMetricModel() *Metrics {
	return new(Metrics)
}

func (m *Metrics) DecodeMetricRequest(body io.ReadCloser) error {
	dec := json.NewDecoder(body)
	if err := dec.Decode(&m); err != nil {
		return err
	}

	if m.MType != GaugeType && m.MType != CounterType {
		return fmt.Errorf("invalid metric type: %s", m.MType)
	}
	if m.MType == GaugeType {
		m.Delta = nil
	} else if m.MType == CounterType {
		m.Value = nil
	}

	return nil
}

type MetricCollection struct {
	Metrics []Metrics `json:"metrics"`
}

func NewMetricCollection() *MetricCollection {
	return &MetricCollection{
		Metrics: []Metrics{},
	}
}

func (mc *MetricCollection) DecodeMetricsRequest(body io.ReadCloser) error {
	var tempMetrics []Metrics

	byteValue, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(byteValue, &tempMetrics); err != nil {
		return err
	}

	for i, m := range tempMetrics {
		if m.MType != GaugeType && m.MType != CounterType {
			return fmt.Errorf("invalid metric type: %s", m.MType)
		}
		if m.MType == GaugeType {
			tempMetrics[i].Delta = nil
		} else if m.MType == CounterType {
			tempMetrics[i].Value = nil
		}
	}

	mc.Metrics = tempMetrics

	return nil
}
