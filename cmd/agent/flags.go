package main

import (
	"flag"
	"time"
)

type ClientConfig interface {
	GetEndpointObject() flag.Value
	GetPollIntervalPointer() *time.Duration
	GetReportIntervalPointer() *time.Duration
	InitAdress()
}

func parseFlags(cfg ClientConfig) {
	flag.Var(cfg.GetEndpointObject(), "a", "net address host:port")
	flag.DurationVar(cfg.GetReportIntervalPointer(), "r", 10*time.Second, "report seconds interval")
	flag.DurationVar(cfg.GetPollIntervalPointer(), "p", 2*time.Second, "poll seconds interval")

	flag.Parse()

	cfg.InitAdress()
}
