package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Address        string `env:"ADDRESS"`
}

func (c config) GetPollInterval() int {
	return c.PollInterval
}

func (c config) GetReportInterval() int {
	return c.ReportInterval
}

func (c config) GetAddress() string {
	return c.Address
}

func (c *config) InitEnv() {
	env.Parse(c)
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.IntVar(&c.ReportInterval, "r", 10, "report seconds interval")
	flag.IntVar(&c.PollInterval, "p", 2, "poll seconds interval")

	flag.Parse()
}

func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()

	return cfg
}
