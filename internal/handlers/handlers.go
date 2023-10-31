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

		if parts[util.METRIC_TYPE] == util.COUNTER_TYPE {
			value, _ := strconv.ParseInt(parts[util.METRIC_VALUE], 10, 64)
			strg.SetCounterMetric(parts[util.METRIC_NAME], storage.Counter(value))
		}

		if parts[util.METRIC_TYPE] == util.GAUGE_TYPE {
			value, _ := strconv.ParseFloat(parts[util.METRIC_VALUE], 64)
			strg.SetGaugeMetric(parts[util.METRIC_NAME], storage.Gauge(value))
		}

		w.WriteHeader(http.StatusOK)
	}
}

// func GetHandle(strg storage.MetricStorage) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("ПРИВЕТ!"))
// 	}
// }
