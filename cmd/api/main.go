package main

import (
	"fmt"
	"net/http"

	"github.com/notOliveira/onde-tem/internal/infra/config"
	"github.com/notOliveira/onde-tem/internal/infra/database"
	"github.com/notOliveira/onde-tem/internal/infra/logger"
)

var (
	log logger.Logger
)

func main() {

	log = *config.GetLogger("main")

	err := config.Init()

	if err != nil {
		log.Errorf("Failed to initialize config: %v", err)
		return
	}

	database.Open(&log)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Onde Tem API 🚀 - Docker up!")
	})
	http.ListenAndServe(":8080", nil)
}
