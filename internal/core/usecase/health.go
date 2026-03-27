package usecase

import (
	"context"
	"github.com/notOliveira/onde-tem/internal/core/ports"
	"time"
)

type HealthUseCase struct {
	cache ports.Cache
}

func NewHealthUseCase(cache ports.Cache) *HealthUseCase {
	return &HealthUseCase{cache: cache}
}

func (u *HealthUseCase) Execute(ctx context.Context) string {
	u.cache.Set(ctx, "teste", "ok", 5*time.Minute)
	val, err := u.cache.Get(ctx, "teste")

	if err == nil && val != "" {
		return "cache hit: " + val
	}

	u.cache.Set(ctx, "teste", "ok", 5*time.Minute)

	return "cache miss"
}
