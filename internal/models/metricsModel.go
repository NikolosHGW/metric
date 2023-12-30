package models

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/NikolosHGW/metric/internal/util"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetricsModel() *Metrics {
	return new(Metrics)
}

func (m *Metrics) DecodeMetricRequest(body io.ReadCloser) error {
	dec := json.NewDecoder(body)
	if err := dec.Decode(&m); err != nil {
		return err
	}

	if m.MType != util.GaugeType && m.MType != util.CounterType {
		return fmt.Errorf("invalid metric type: %s", m.MType)
	}

	return nil
}
