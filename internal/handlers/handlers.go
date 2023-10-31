package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/NikolosHGW/metric/internal/storage"
	"github.com/NikolosHGW/metric/internal/util"
)

func PostHandle(strg storage.MetricStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := util.SliceStrings(strings.Split(r.URL.Path, "/"), 0)

		if parts[util.MetricType] == util.CounterType {
			value, _ := strconv.ParseInt(parts[util.MetricValue], 10, 64)
			strg.SetCounterMetric(parts[util.MetricName], storage.Counter(value))
		}

		if parts[util.MetricType] == util.GaugeType {
			value, _ := strconv.ParseFloat(parts[util.MetricValue], 64)
			strg.SetGaugeMetric(parts[util.MetricName], storage.Gauge(value))
		}

		w.WriteHeader(http.StatusOK)
	}
}
