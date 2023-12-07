package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/NikolosHGW/metric/internal/util"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Sugar().Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}

	return http.HandlerFunc(logFn)
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

		isCounterType := util.CheckCounterType(parts[util.MetricType], parts[util.MetricValue])
		isGaugeType := util.CheckGaugeType(parts[util.MetricType], parts[util.MetricValue])

		if !isCounterType && !isGaugeType {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}
