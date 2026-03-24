package config

import "github.com/notOliveira/onde-tem/internal/infra/logger"

func GetLogger(p string) *logger.Logger {
	return logger.NewLogger(p)
}

func Init() error {
	log := GetLogger("config")

	log.Info("Config initialized")
	return nil
}
