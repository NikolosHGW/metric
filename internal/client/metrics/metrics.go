package metrics

import (
	"math/rand"
	"runtime"

	"github.com/NikolosHGW/metric/internal/models"
)

type Metrics struct {
	runtime.MemStats
	PollCount   models.Counter
	RandomValue models.Gauge
}

func (m *Metrics) IncPollCount() {
	m.PollCount += 1
}

func (m *Metrics) UpdateRandomValue() {
	m.RandomValue = models.Gauge(rand.Float64())
}

func (m *Metrics) RefreshMetrics() {
	runtime.ReadMemStats(&m.MemStats)
}

func (m Metrics) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		models.Alloc:         models.Gauge(m.Alloc),
		models.BuckHashSys:   models.Gauge(m.BuckHashSys),
		models.Frees:         models.Gauge(m.Frees),
		models.GCCPUFraction: models.Gauge(m.GCCPUFraction),
		models.GCSys:         models.Gauge(m.GCSys),
		models.HeapAlloc:     models.Gauge(m.HeapAlloc),
		models.HeapIdle:      models.Gauge(m.HeapIdle),
		models.HeapInuse:     models.Gauge(m.HeapInuse),
		models.HeapObjects:   models.Gauge(m.HeapObjects),
		models.HeapReleased:  models.Gauge(m.HeapReleased),
		models.HeapSys:       models.Gauge(m.HeapSys),
		models.LastGC:        models.Gauge(m.LastGC),
		models.Lookups:       models.Gauge(m.Lookups),
		models.MCacheInuse:   models.Gauge(m.MCacheInuse),
		models.MCacheSys:     models.Gauge(m.MCacheSys),
		models.MSpanInuse:    models.Gauge(m.MSpanInuse),
		models.MSpanSys:      models.Gauge(m.MSpanSys),
		models.Mallocs:       models.Gauge(m.Mallocs),
		models.NextGC:        models.Gauge(m.NextGC),
		models.NumForcedGC:   models.Gauge(m.NumForcedGC),
		models.NumGC:         models.Gauge(m.NumGC),
		models.OtherSys:      models.Gauge(m.OtherSys),
		models.PauseTotalNs:  models.Gauge(m.PauseTotalNs),
		models.StackInuse:    models.Gauge(m.StackInuse),
		models.StackSys:      models.Gauge(m.StackSys),
		models.Sys:           models.Gauge(m.Sys),
		models.TotalAlloc:    models.Gauge(m.TotalAlloc),
		models.PollCount:     models.Counter(m.PollCount),
		models.RandomValue:   models.Gauge(m.RandomValue),
	}
}

func NewMetrics() *Metrics {
	return new(Metrics)
}
