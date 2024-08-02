package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
)

type config struct {
	Address        string `env:"ADDRESS" json:"address,omitempty"`
	Key            string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	ConfigPath     string `env:"CONFIG"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval,omitempty"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval,omitempty"`
	RateLimit      int    `end:"RATE_LIMIT"`
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

func (c config) GetKey() string {
	return c.Key
}

func (c config) GetRateLimit() int {
	return c.RateLimit
}

func (c *config) InitEnv() {
	err := env.Parse(c)
	if err != nil {
		log.Println("err when parse env")
	}
}

func (c config) GetCryptoKeyPath() string {
	return c.CryptoKey
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.IntVar(&c.ReportInterval, "r", 10, "report seconds interval")
	flag.IntVar(&c.PollInterval, "p", 2, "poll seconds interval")
	flag.StringVar(&c.Key, "k", "", "secret key for hash")
	flag.IntVar(&c.RateLimit, "l", 10, "Rate limit for outgoing requests")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "path to public crypto key")
	flag.StringVar(&c.ConfigPath, "c", "", "path to config file")

	flag.Parse()
}

func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()
	cfg.loadFromJSON()

	return cfg
}

func (c *config) loadFromJSON() {
	if c.ConfigPath == "" {
		return
	}

	fileContent, err := os.ReadFile(c.ConfigPath)
	if err != nil {
		log.Println("could not read config file: %w", err)
	}

	tempConfig := config{}
	if err = json.Unmarshal(fileContent, &tempConfig); err != nil {
		log.Println("invalid config file content: %w", err)
	}

	if c.Address == "" && tempConfig.Address != "" {
		c.Address = tempConfig.Address
	}

	if c.CryptoKey == "" && tempConfig.CryptoKey != "" {
		c.CryptoKey = tempConfig.CryptoKey
	}

	if c.PollInterval == 2 && tempConfig.PollInterval != 2 {
		c.PollInterval = tempConfig.PollInterval
	}

	if c.ReportInterval == 10 && tempConfig.ReportInterval != 10 {
		c.ReportInterval = tempConfig.ReportInterval
	}
}
