package usecase

import (
	"context"

	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/ports"
)

type ListEstablishmentsUseCase struct {
	repo ports.EstablishmentRepository
}

func NewListEstablishmentsUseCase(repo ports.EstablishmentRepository) *ListEstablishmentsUseCase {
	return &ListEstablishmentsUseCase{
		repo: repo,
	}
}

func (uc *ListEstablishmentsUseCase) Execute(ctx context.Context, filter *domain.EstablishmentFilter) ([]*domain.Establishment, error) {
	return uc.repo.List(ctx, filter)
}
