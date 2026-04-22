package api

import (
	"github.com/gin-gonic/gin"
	inboundHttp "github.com/notOliveira/onde-tem/internal/adapters/inbound/http"
)

type RouterConfig struct {
	EstablishmentHandler *inboundHttp.EstablishmentHandler
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		establishments := v1.Group("/establishments")
		{
			establishments.POST("", cfg.EstablishmentHandler.HandleCreate)
		}
	}

	return router
}
