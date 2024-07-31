package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type config struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DBConnect       string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
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
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "path where store metrics")
	flag.BoolVar(&c.Restore, "r", true, "need load from file")
	flag.StringVar(&c.DBConnect, "d", "user=nikolos password=abc123 dbname=metric sslmode=disable", "data source name for connection")
	flag.StringVar(&c.Key, "k", "", "secret key for hash")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "path to private crypto key")
	flag.Parse()
}

// NewConfig конструктор конфига, в котором идёт инициализация флагов и env переменных
func NewConfig() *config {
	cfg := new(config)

	cfg.parseFlags()
	cfg.InitEnv()

	return cfg
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
