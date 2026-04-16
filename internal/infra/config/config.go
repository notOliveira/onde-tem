package config

import (
	"github.com/notOliveira/onde-tem/internal/infra/logger"
	"os"
)

type Config struct {
	ServerPort string
	ValkeyAddr string
}

func LoadConfig() *Config {
	host := os.Getenv("VALKEY_HOST")
	port := os.Getenv("VALKEY_PORT")

	server := os.Getenv("SERVER_PORT")

	if server != "" && server[0] != ':' {
		server = ":" + server
	}

	return &Config{
		ServerPort: server,
		ValkeyAddr: host + ":" + port,
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
