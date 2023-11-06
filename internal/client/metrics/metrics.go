package metrics

import (
	"runtime"

	"github.com/NikolosHGW/metric/internal/util"
)

type ClientMetrics interface {
	IncPollCount()
	UpdateRandomValue()
	RefreshMetrics()
	GetMetrics() map[string]interface{}
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

func (m Metrics) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		util.Alloc:         m.Alloc,
		util.BuckHashSys:   m.BuckHashSys,
		util.Frees:         m.Frees,
		util.GCCPUFraction: m.GCCPUFraction,
		util.GCSys:         m.GCSys,
		util.HeapAlloc:     m.HeapAlloc,
		util.HeapIdle:      m.HeapIdle,
		util.HeapInuse:     m.HeapInuse,
		util.HeapObjects:   m.HeapObjects,
		util.HeapReleased:  m.HeapReleased,
		util.HeapSys:       m.HeapSys,
		util.LastGC:        m.LastGC,
		util.Lookups:       m.Lookups,
		util.MCacheInuse:   m.MCacheInuse,
		util.MCacheSys:     m.MCacheSys,
		util.MSpanInuse:    m.MSpanInuse,
		util.MSpanSys:      m.MSpanSys,
		util.Mallocs:       m.Mallocs,
		util.NextGC:        m.NextGC,
		util.NumForcedGC:   m.NumForcedGC,
		util.NumGC:         m.NumGC,
		util.OtherSys:      m.OtherSys,
		util.PauseTotalNs:  m.PauseTotalNs,
		util.StackInuse:    m.StackInuse,
		util.StackSys:      m.StackSys,
		util.Sys:           m.Sys,
		util.TotalAlloc:    m.TotalAlloc,
		util.PollCount:     m.PollCount,
		util.RandomValue:   m.RandomValue,
	}
}

func NewMetrics() *Metrics {
	return new(Metrics)
}
