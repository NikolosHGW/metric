package main

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/NikolosHGW/metric/internal/server/util"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	strg := storage.InitStorage()

	mux.Handle(
		"/update/",
		util.MiddlewareConveyor(
			http.HandlerFunc(handlers.PostHandle(strg)),
			middlewares.CheckTypeAndValueMiddleware,
			middlewares.CheckMetricNameMiddleware,
			middlewares.CheckPostMiddleware,
		),
	)

	return http.ListenAndServe(":8080", mux)
}
