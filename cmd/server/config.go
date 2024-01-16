package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	storeInterval   int    `env:"STORE_INTERVAL"`
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
	restore         bool   `env:"RESTORE"`
}

func (c *config) InitEnv() {
	env.Parse(c)
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.storeInterval, "i", 300, "store metrics to file seconds interval")
	flag.StringVar(&c.fileStoragePath, "f", "/tmp/metrics-db.json", "path where store metrics")
	flag.BoolVar(&c.restore, "r", true, "need load from file")
	flag.Parse()
}

func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()

	return cfg
}

func (c config) GetStoreInterval() int {
	return c.storeInterval
}

func (c config) GetFileStoragePath() string {
	return c.fileStoragePath
}

func (c config) GetRestore() bool {
	return c.restore
}
