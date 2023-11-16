package routes

import (
	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/routes/update"
	"github.com/NikolosHGW/metric/internal/server/routes/value"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/go-chi/chi"
)

func InitRouter(strg storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.WithSetMetricHandle((strg)))

		update.InitUpdateRoutes(r, strg)
		value.InitValueRoutes(r, strg)
	})

	return r
}
