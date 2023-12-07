package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

func (c *config) InitEnv() {
	env.Parse(c)
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.Parse()
}

func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()

	return cfg
}
