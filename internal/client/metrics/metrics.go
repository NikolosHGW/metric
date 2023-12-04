package metrics

import (
	"math/rand"
	"runtime"

	"github.com/NikolosHGW/metric/internal/util"
)

type Metrics struct {
	runtime.MemStats
	PollCount   util.Counter
	RandomValue util.Gauge
}

func (m *Metrics) IncPollCount() {
	m.PollCount += 1
}

func (m *Metrics) UpdateRandomValue() {
	m.RandomValue = util.Gauge(rand.Float64())
}

func (m *Metrics) RefreshMetrics() {
	runtime.ReadMemStats(&m.MemStats)
}

func (m Metrics) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		util.Alloc:         util.Gauge(m.Alloc),
		util.BuckHashSys:   util.Gauge(m.BuckHashSys),
		util.Frees:         util.Gauge(m.Frees),
		util.GCCPUFraction: util.Gauge(m.GCCPUFraction),
		util.GCSys:         util.Gauge(m.GCSys),
		util.HeapAlloc:     util.Gauge(m.HeapAlloc),
		util.HeapIdle:      util.Gauge(m.HeapIdle),
		util.HeapInuse:     util.Gauge(m.HeapInuse),
		util.HeapObjects:   util.Gauge(m.HeapObjects),
		util.HeapReleased:  util.Gauge(m.HeapReleased),
		util.HeapSys:       util.Gauge(m.HeapSys),
		util.LastGC:        util.Gauge(m.LastGC),
		util.Lookups:       util.Gauge(m.Lookups),
		util.MCacheInuse:   util.Gauge(m.MCacheInuse),
		util.MCacheSys:     util.Gauge(m.MCacheSys),
		util.MSpanInuse:    util.Gauge(m.MSpanInuse),
		util.MSpanSys:      util.Gauge(m.MSpanSys),
		util.Mallocs:       util.Gauge(m.Mallocs),
		util.NextGC:        util.Gauge(m.NextGC),
		util.NumForcedGC:   util.Gauge(m.NumForcedGC),
		util.NumGC:         util.Gauge(m.NumGC),
		util.OtherSys:      util.Gauge(m.OtherSys),
		util.PauseTotalNs:  util.Gauge(m.PauseTotalNs),
		util.StackInuse:    util.Gauge(m.StackInuse),
		util.StackSys:      util.Gauge(m.StackSys),
		util.Sys:           util.Gauge(m.Sys),
		util.TotalAlloc:    util.Gauge(m.TotalAlloc),
		util.PollCount:     util.Counter(m.PollCount),
		util.RandomValue:   util.Gauge(m.RandomValue),
	}
}

func NewMetrics() *Metrics {
	return new(Metrics)
}
