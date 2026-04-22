package main

import (
	inboundHttp "github.com/notOliveira/onde-tem/internal/adapters/inbound/http"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/cache"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres/sqlc"
	"github.com/notOliveira/onde-tem/internal/core/usecase"
	"github.com/notOliveira/onde-tem/internal/infra/api"
	"github.com/notOliveira/onde-tem/internal/infra/config"
	"github.com/notOliveira/onde-tem/internal/infra/database"
)

func main() {
	log := config.GetLogger("main")
	cfg := config.LoadConfig()

	connPool, err := database.NewConnection(log)
	if err != nil {
		log.Errorf("failed to connect db: %v", err)
	}
	defer connPool.Close()

	queries := sqlc.New(connPool)
	establishmentRepo := postgres.NewEstablishmentRepository(queries)

	valkeyClient := cache.NewValkeyClient(cfg.ValkeyAddr)
	_ = valkeyClient

	createEstablishmentUC := usecase.NewCreateEstablishmentUseCase(establishmentRepo)

	establishmentHandler := inboundHttp.NewEstablishmentHandler(createEstablishmentUC)

	router := api.NewRouter(api.RouterConfig{
		EstablishmentHandler: establishmentHandler,
	})

	log.Infof("Server running on %s", cfg.ServerPort)

	if err := router.Run(cfg.ServerPort); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
