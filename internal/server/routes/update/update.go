package update

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func InitUpdateRoutes(r chi.Router, h http.HandlerFunc) {
	r.Route("/update", func(r chi.Router) {
		r.Use(middlewares.CheckMetricNameMiddleware)
		r.Use(middlewares.CheckTypeAndValueMiddleware)

		r.Post("/{metricType}/{metricName}/{metricValue}", h)
	})
}
