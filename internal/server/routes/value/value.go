package value

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitValueRoutes(r chi.Router, h http.HandlerFunc, hJson http.HandlerFunc) {
	r.Route("/value", func(r chi.Router) {
		r.Post("/", hJson)
		r.Get("/{metricType}/{metricName}", h)
	})
}
