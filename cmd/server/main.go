package main

import (
	"net/http"

	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/storage/memory"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	strg := memory.NewMemStorage()

	r := routes.InitRouter(strg)

	return http.ListenAndServe(":8080", r)
}
