package routes

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/routes/update"
	"github.com/NikolosHGW/metric/internal/server/routes/value"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler interface {
	SetMetric(http.ResponseWriter, *http.Request)
	SetJSONMetric(http.ResponseWriter, *http.Request)
	GetMetric(http.ResponseWriter, *http.Request)
	GetValueMetric(http.ResponseWriter, *http.Request)
	GetMetrics(http.ResponseWriter, *http.Request)
	PingDB(http.ResponseWriter, *http.Request)
	UpsertMetrics(http.ResponseWriter, *http.Request)
}

type Middleware interface {
	WithHash(http.Handler) http.Handler
}

type DecryptMiddleware interface {
	DecryptHandler(next http.Handler) http.Handler
}

type CheckIPMiddleware interface {
	WithCheckIP(next http.Handler) http.Handler
}

func InitRouter(handler Handler, myMiddleware Middleware, decryptMiddleware DecryptMiddleware, checkIP CheckIPMiddleware) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.WithLogging)
	r.Use(middlewares.WithGzip)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetMetrics)
		r.Get("/ping", handler.PingDB)
		r.With(myMiddleware.WithHash, decryptMiddleware.DecryptHandler, checkIP.WithCheckIP).Post("/updates/", handler.UpsertMetrics)

		update.InitUpdateRoutes(r, handler.SetMetric, handler.SetJSONMetric)
		value.InitValueRoutes(r, handler.GetValueMetric, handler.GetMetric)
	})

	r.Mount("/debug", middleware.Profiler())

	return r
}
