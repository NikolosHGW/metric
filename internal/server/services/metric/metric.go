package metric

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/NikolosHGW/metric/internal/util"
)

func SetMetric(r *http.Request, strg storage.Storage) {
	parts := util.SliceStrings(strings.Split(r.URL.Path, "/"), 0)

	if parts[util.MetricType] == util.CounterType {
		value, _ := strconv.ParseInt(parts[util.MetricValue], 10, 64)
		strg.SetCounterMetric(parts[util.MetricName], util.Counter(value))
	}

	if parts[util.MetricType] == util.GaugeType {
		value, _ := strconv.ParseFloat(parts[util.MetricValue], 64)
		strg.SetGaugeMetric(parts[util.MetricName], util.Gauge(value))
	}
}

func GetMetricValue() string {
	return "12"
}
