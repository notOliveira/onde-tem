package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/notOliveira/onde-tem/internal/adapters/outbound/cache"
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

	cacheClient := cache.NewValkeyClient("valkey:6379")

	ctx := context.Background()

	err = cacheClient.Set(ctx, "teste", "ok", 0)
	if err != nil {
		log.Errorf("cache error: %v", err)
	}

	val, _ := cacheClient.Get(ctx, "teste")
	log.Infof("cache value: %s", val)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Onde Tem API 🚀 com cache! Variável em cache: %s", val)
	})

	log.Info("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}
