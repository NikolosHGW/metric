package main

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/handlers"
	"github.com/NikolosHGW/metric/internal/middlewares"
	"github.com/NikolosHGW/metric/internal/util"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	mux.Handle("/update/", util.MiddlewareConveyor(http.HandlerFunc(handlers.PostHandle), middlewares.CheckPostMiddleware))

	return http.ListenAndServe(":8080", mux)
}
