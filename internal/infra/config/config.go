package config

import (
	"os"
	"strconv"

	"github.com/notOliveira/onde-tem/internal/infra/logger"
)

type Config struct {
	ServerPort string
	ValkeyAddr string
	CacheTTL   int
}

func LoadConfig() *Config {
	host := os.Getenv("VALKEY_HOST")
	port := os.Getenv("VALKEY_PORT")
	server := os.Getenv("SERVER_PORT")
	if server != "" && server[0] != ':' {
		server = ":" + server
	}
	cacheTTL := os.Getenv("CACHE_TTL_SECONDS")
	cacheTTLInt, err := strconv.Atoi(cacheTTL)
	if err != nil || cacheTTLInt <= 0 {
		cacheTTLInt = 600
	}

	return &Config{
		ServerPort: server,
		ValkeyAddr: host + ":" + port,
		CacheTTL:   cacheTTLInt,
	}
}

func GetLogger(p string) *logger.Logger {
	return logger.NewLogger(p)
}

func Init() error {
	log := GetLogger("config")

	log.Info("Config initialized")
	return nil
}
