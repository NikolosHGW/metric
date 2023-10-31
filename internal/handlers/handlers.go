package handlers

import "net/http"

func PostHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello WORLD!!!"))
}
