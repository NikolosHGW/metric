package routes

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/routes/update"
	"github.com/NikolosHGW/metric/internal/server/routes/value"
	"github.com/go-chi/chi"
)

type Handler interface {
	SetMetric(http.ResponseWriter, *http.Request)
	GetValueMetric(http.ResponseWriter, *http.Request)
	GetMetrics(http.ResponseWriter, *http.Request)
}

func InitRouter(handler Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.WithLogging)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetMetrics)

		update.InitUpdateRoutes(r, handler.SetMetric)
		value.InitValueRoutes(r, handler.GetValueMetric)
	})

	return r
}
