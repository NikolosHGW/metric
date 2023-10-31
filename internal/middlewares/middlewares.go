package middlewares

import (
	"net/http"
	"strings"

	"github.com/NikolosHGW/metric/internal/util"
)

func CheckPostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckMetricNameMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := util.SliceStrings(strings.Split(r.URL.Path, "/"), 0)

		if len(parts) != 4 {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckTypeAndValueMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := util.SliceStrings(strings.Split(r.URL.Path, "/"), 0)

		isCounterType := util.CheckCounterType(parts[util.METRIC_TYPE], parts[util.METRIC_VALUE])
		isGaugeType := util.CheckGaugeType(parts[util.METRIC_TYPE], parts[util.METRIC_VALUE])

		if !isCounterType && !isGaugeType {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}
