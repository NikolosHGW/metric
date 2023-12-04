package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	Address string `env:"ADDRESS"`
}

func (c *config) InitEnv() {
	env.Parse(c)
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.Parse()
}

func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()

	return cfg
}
