package update

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitUpdateRoutes(r chi.Router, h http.HandlerFunc) {
	r.Post("/update", h)
}
