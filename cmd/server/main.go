package main

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/storage/memory"
	"github.com/go-chi/chi"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	strg := memory.NewMemStorage()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.PostHandle((strg)))
		r.Route("/update", func(r chi.Router) {
			r.Use(middlewares.CheckMetricNameMiddleware)
			r.Use(middlewares.CheckTypeAndValueMiddleware)

			r.Post("/{metricType}/{metricName}/{metricValue}", handlers.PostHandle(strg))
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/{metricType}/{metricName}", handlers.PostHandle(strg))
		})
	})

	return http.ListenAndServe(":8080", r)
}
