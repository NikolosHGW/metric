package update

import (
	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/go-chi/chi"
)

func InitUpdateRoutes(r chi.Router, strg storage.Storage) {
	r.Route("/update", func(r chi.Router) {
		r.Use(middlewares.CheckMetricNameMiddleware)
		r.Use(middlewares.CheckTypeAndValueMiddleware)

		r.Post("/{metricType}/{metricName}/{metricValue}", handlers.PostHandle(strg))
	})
}
