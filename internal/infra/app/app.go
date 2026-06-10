package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	inboundHttp "github.com/notOliveira/onde-tem/internal/adapters/inbound/http"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/cache"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres/sqlc"
	"github.com/notOliveira/onde-tem/internal/core/usecase"
	"github.com/notOliveira/onde-tem/internal/infra/api"
	"github.com/notOliveira/onde-tem/internal/infra/config"
	"github.com/notOliveira/onde-tem/internal/infra/database"
)

type App struct {
	Router *gin.Engine
	Config *config.Config
	DB     *pgxpool.Pool
}

func NewApp() (*App, error) {
	cfg := config.LoadConfig()
	log := config.GetLogger("main")

	connPool, err := database.NewConnection(log)
	if err != nil {
		return nil, err
	}

	// Initialize outbound adapters
	valkeyClient := cache.NewValkeyClient(cfg.ValkeyAddr)

	// Initialize repositories
	queries := sqlc.New(connPool)
	repo := postgres.NewEstablishmentRepository(queries)
	cachedRepo := cache.NewCachedRepository(repo, valkeyClient, time.Duration(cfg.CacheTTL)*time.Second)

	// Initialize use cases
	createUC := usecase.NewCreateEstablishmentUseCase(cachedRepo)
	listUC := usecase.NewListEstablishmentsUseCase(cachedRepo)

	// Initialize handlers
	handler := inboundHttp.NewEstablishmentHandler(
		createUC,
		listUC,
	)

	// Initialize router
	router := api.NewRouter(api.RouterConfig{
		EstablishmentHandler: handler,
	})

	return &App{
		Router: router,
		Config: cfg,
		DB:     connPool,
	}, nil
}

func (a *App) Close() {
	a.DB.Close()
}
