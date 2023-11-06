package metrics

import (
	"runtime"

	"github.com/NikolosHGW/metric/internal/util"
)

type ClientMetrics interface {
	IncPollCount()
	UpdateRandomValue()
	RefreshMetrics()
	GetMetrics() []interface{}
}

type Metrics struct {
	runtime.MemStats
	PollCount   util.Counter
	RandomValue util.Gauge
}

func (m *Metrics) IncPollCount() {
	m.PollCount += 1
}

func (m *Metrics) UpdateRandomValue() {
	m.RandomValue = 901.02
}

func (m *Metrics) RefreshMetrics() {
	runtime.ReadMemStats(&m.MemStats)
}

func (m Metrics) GetMetrics() []interface{} {
	return []interface{}{m.Alloc, m.PollCount, m.RandomValue}
}

func NewMetrics() *Metrics {
	return new(Metrics)
}
