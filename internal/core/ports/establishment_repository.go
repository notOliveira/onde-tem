package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/notOliveira/onde-tem/internal/core/domain"
)

type EstablishmentRepository interface {
	Create(ctx context.Context, establishment *domain.Establishment) error
	Update(ctx context.Context, establishment *domain.Establishment) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Establishment, error)
	GetBySlug(ctx context.Context, slug domain.Slug) (*domain.Establishment, error)
	List(ctx context.Context, filter *domain.EstablishmentFilter) ([]*domain.Establishment, error)
}
