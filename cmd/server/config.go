package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DBConnect       string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}

func (c *config) InitEnv() {
	env.Parse(c)
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.StoreInterval, "i", 300, "store metrics to file seconds interval")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "path where store metrics")
	flag.BoolVar(&c.Restore, "r", true, "need load from file")
	flag.StringVar(&c.DBConnect, "d", "user=nikolos password=abc123 dbname=metric sslmode=disable", "data source name for connection")
	flag.StringVar(&c.Key, "k", "", "secret key for hash")
	flag.Parse()
}

func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()

	return cfg
}

func (c config) GetStoreInterval() int {
	return c.StoreInterval
}

func (c config) GetFileStoragePath() string {
	return c.FileStoragePath
}

func (c config) GetRestore() bool {
	return c.Restore
}

func (c config) GetDBConnection() string {
	return c.DBConnect
}

func (c config) GetKey() string {
	return c.Key
}
