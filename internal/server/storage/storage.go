package storage

import (
	"github.com/NikolosHGW/metric/internal/util"
)

type Storage interface {
	SetGaugeMetric(string, util.Gauge)
	SetCounterMetric(string, util.Counter)
	GetGaugeMetric(string) (util.Gauge, error)
	GetCounterMetric(string) (util.Counter, error)
}
