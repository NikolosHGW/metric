package value

import (
	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/go-chi/chi"
)

func InitValueRoutes(r chi.Router, strg storage.Storage) {
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", handlers.PostHandle(strg))
	})
}
