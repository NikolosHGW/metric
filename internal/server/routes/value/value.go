package value

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitValueRoutes(r chi.Router, h http.HandlerFunc) {
	r.Post("/value", h)
}
