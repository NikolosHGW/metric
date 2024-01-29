package update

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func InitUpdateRoutes(r chi.Router, h http.HandlerFunc, hJSON http.HandlerFunc) {
	r.Route("/update", func(r chi.Router) {
		r.Post("/", hJSON)

		textGroup := r.Group(nil)
		textGroup.Use(middlewares.CheckMetricNameMiddleware)
		textGroup.Use(middlewares.CheckTypeAndValueMiddleware)
		textGroup.Post("/{metricType}/{metricName}/{metricValue}", h)
	})
}
