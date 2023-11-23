package main

import (
	"flag"
)

type ClientConfig interface {
	GetEndpointObject() flag.Value
	GetPollIntervalPointer() *int
	GetReportIntervalPointer() *int
	InitAdress()
}

func parseFlags(cfg ClientConfig) {
	flag.Var(cfg.GetEndpointObject(), "a", "net address host:port")
	flag.IntVar(cfg.GetReportIntervalPointer(), "r", 10, "report seconds interval")
	flag.IntVar(cfg.GetPollIntervalPointer(), "p", 2, "poll seconds interval")

	flag.Parse()

	cfg.InitAdress()
}
