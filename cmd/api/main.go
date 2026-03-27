package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/notOliveira/onde-tem/internal/adapters/outbound/cache"
	"github.com/notOliveira/onde-tem/internal/core/usecase"
	"github.com/notOliveira/onde-tem/internal/infra/config"
	"github.com/notOliveira/onde-tem/internal/infra/database"
)

func main() {

	log := config.GetLogger("main")

	if err := config.Init(); err != nil {
		log.Errorf("Failed to initialize config: %v", err)
		return
	}

	conn, err := database.NewConnection(log)
	if err != nil {
		log.Errorf("failed to connect db: %v", err)
		return
	}
	defer conn.Close(context.Background())

	host := os.Getenv("VALKEY_HOST")
	port := os.Getenv("VALKEY_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	cacheClient := cache.NewValkeyClient(addr)

	// USECASE
	healthUC := usecase.NewHealthUseCase(cacheClient)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := healthUC.Execute(r.Context())
		fmt.Fprintln(w, result)
	})

	log.Info("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}
