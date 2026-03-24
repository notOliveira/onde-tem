package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/notOliveira/onde-tem/internal/infra/logger"
)

var DB *pgx.Conn

func Open(log *logger.Logger) {

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error

	DB, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Errorf("Unable to connect to database: %v", err)
		os.Exit(1)
	}

	log.Infof("Database connected")

	err = DB.Ping(context.Background())
	if err != nil {
		log.Errorf("Database ping failed: %v", err)
		os.Exit(1)
	} else {
		log.Infof("Database ping successful")
	}
}

func Close() {
	if DB != nil {
		err := DB.Close(context.Background())
		if err != nil {
			fmt.Printf("Error closing database connection: %v\n", err)
		} else {
			fmt.Println("Database connection closed")
		}
	}
}
