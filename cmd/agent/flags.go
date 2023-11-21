package main

import (
	"flag"

	"github.com/NikolosHGW/metric/internal/client/config"
)

func parseFlags(configFlags *config.Flags) {
	flag.Var(&configFlags.Endpoint, "a", "net address host:port")
	flag.IntVar(&configFlags.ReportInterval, "r", 10, "report seconds interval")
	flag.IntVar(&configFlags.PollInterval, "p", 2, "poll seconds interval")

	flag.Parse()
}
