package usecase

import (
	"context"

	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/ports"
)

type CreateEstablishmentInput struct {
	Name     string
	Slug     string
	Types    []string
	Location domain.Location
	Address  domain.Address
}

type CreateEstablishmentUseCase struct {
	repo  ports.EstablishmentRepository
	cache ports.Cache
}

func NewCreateEstablishmentUseCase(repo ports.EstablishmentRepository, cache ports.Cache) *CreateEstablishmentUseCase {
	return &CreateEstablishmentUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *CreateEstablishmentUseCase) Execute(ctx context.Context, input CreateEstablishmentInput) (*domain.Establishment, error) {

	var types []domain.EstablishmentType
	for _, t := range input.Types {
		types = append(types, domain.EstablishmentType(t))
	}

	est, err := domain.NewEstablishment(
		input.Name,
		input.Slug,
		types,
		input.Location,
		input.Address,
	)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(ctx, est)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.Delete(ctx, "est-"+est.ID().String()); err != nil {
		return nil, err
	}

	return est, nil
}
