package value

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitValueRoutes(r chi.Router, h http.HandlerFunc, hJSON http.HandlerFunc) {
	r.Route("/value", func(r chi.Router) {
		r.Post("/", hJSON)
		r.Get("/{metricType}/{metricName}", h)
	})
}
