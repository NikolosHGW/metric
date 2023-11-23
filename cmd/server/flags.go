package main

import "flag"

type ServerConfig interface {
	GetEndpointObject() flag.Value
	InitAdress()
}

func parseFlags(cfg ServerConfig) {
	flag.Var(cfg.GetEndpointObject(), "a", "net address host:port")
	flag.Parse()

	cfg.InitAdress()
}
