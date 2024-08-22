package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
)

const (
	DefaultFileStoragePath = "/tmp/metrics-db.json"
	DefaultDBConnect       = "user=nikolos password=abc123 dbname=metric sslmode=disable"
)

type config struct {
	Address         string `env:"ADDRESS" json:"address,omitempty"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file,omitempty"`
	DBConnect       string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	Key             string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	ConfigPath      string `env:"CONFIG"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet,omitempty"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval,omitempty"`
	Restore         bool   `env:"RESTORE" json:"restore,omitempty"`
}

func (c *config) InitEnv() {
	err := env.Parse(c)
	if err != nil {
		log.Println("err when parse env")
	}
}

func (c *config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.StoreInterval, "i", 300, "store metrics to file seconds interval")
	flag.StringVar(&c.FileStoragePath, "f", DefaultFileStoragePath, "path where store metrics")
	flag.BoolVar(&c.Restore, "r", true, "need load from file")
	flag.StringVar(&c.DBConnect, "d", DefaultDBConnect, "data source name for connection")
	flag.StringVar(&c.Key, "k", "", "secret key for hash")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "path to private crypto key")
	flag.StringVar(&c.ConfigPath, "c", "", "path to config file")
	flag.StringVar(&c.TrustedSubnet, "t", "", "trusted subnet CIDR")
	flag.Parse()
}

// NewConfig конструктор конфига, в котором идёт инициализация флагов и env переменных
func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()
	cfg.loadFromJSON()

	return cfg
}

// GetAddress геттер для хоста
func (c config) GetAddress() string {
	return c.Address
}

// GetStoreInterval геттер для интервала хранения
func (c config) GetStoreInterval() int {
	return c.StoreInterval
}

// GetFileStoragePath геттер для пути к хранению
func (c config) GetFileStoragePath() string {
	return c.FileStoragePath
}

// GetRestore геттер для флага нужно ли хранить метрики на диске
func (c config) GetRestore() bool {
	return c.Restore
}

// GetDBConnection геттер для подключения к бд
func (c config) GetDBConnection() string {
	return c.DBConnect
}

// GetKey геттер для секретного ключа для хеширования
func (c config) GetKey() string {
	return c.Key
}

// GetCryptoKeyPath геттер для пути к приватному ключу шифрования
func (c config) GetCryptoKeyPath() string {
	return c.CryptoKey
}

// GetTrustedSubnet геттер для CIDR
func (c config) GetTrustedSubnet() string {
	return c.TrustedSubnet
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

	if c.Restore && !tempConfig.Restore {
		c.Restore = tempConfig.Restore
	}

	if c.StoreInterval == 300 && tempConfig.StoreInterval != 300 {
		c.StoreInterval = tempConfig.StoreInterval
	}

	if c.FileStoragePath == DefaultFileStoragePath && tempConfig.FileStoragePath != DefaultFileStoragePath {
		c.FileStoragePath = tempConfig.FileStoragePath
	}

	if c.DBConnect == DefaultDBConnect && tempConfig.FileStoragePath != DefaultDBConnect {
		c.FileStoragePath = tempConfig.FileStoragePath
	}

	if c.CryptoKey == "" && tempConfig.CryptoKey != "" {
		c.CryptoKey = tempConfig.CryptoKey
	}

	if c.TrustedSubnet == "" && tempConfig.TrustedSubnet != "" {
		c.TrustedSubnet = tempConfig.TrustedSubnet
	}
}
