package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/notOliveira/onde-tem/internal/infra/logger"
)

func NewConnection(log *logger.Logger) (*pgx.Conn, error) {

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Errorf("Unable to connect to database: %v", err)
		return nil, err
	}

	log.Infof("Database connected")

	if err := conn.Ping(context.Background()); err != nil {
		log.Errorf("Database ping failed: %v", err)
		return nil, err
	}

	log.Infof("Database ping successful")

	return conn, nil
}
