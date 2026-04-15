package config

import (
	"os"
	"github.com/notOliveira/onde-tem/internal/infra/logger"
)


type Config struct {
    ServerPort string
    ValkeyAddr string
}

func LoadConfig() *Config {
    host := os.Getenv("VALKEY_HOST")
    port := os.Getenv("VALKEY_PORT")
    
    return &Config{
        ServerPort: os.Getenv("SERVER_PORT"),
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
