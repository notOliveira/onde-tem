package postgres

import (
	"context"

	"errors"
	"github.com/jackc/pgx/v5"

	"encoding/json"
	"github.com/google/uuid"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres/sqlc"
	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/ports"
)

type establishmentRepository struct {
	q *sqlc.Queries
}

func NewEstablishmentRepository(q *sqlc.Queries) ports.EstablishmentRepository {
	return &establishmentRepository{q: q}
}

func (r *establishmentRepository) Create(
	ctx context.Context,
	e *domain.Establishment,
) error {

	phonesJSON, err := json.Marshal(e.Phones())
	if err != nil {
		return err
	}

	addressJSON, err := json.Marshal(e.Address())
	if err != nil {
		return err
	}

	types := make([]string, len(e.Types()))
	for i, t := range e.Types() {
		types[i] = t.String()
	}

	params := sqlc.CreateEstablishmentParams{
		ID:        e.ID(),
		Name:      e.Name(),
		Slug:      e.Slug().String(),
		Types:     types,
		Email:     e.Email(),
		Website:   e.Website(),
		Timezone:  e.Timezone(),
		Phones:    phonesJSON,
		Address:   addressJSON,
		Lat:       e.Location().Lat(),
		Lon:       e.Location().Lon(),
		CreatedAt: e.CreatedAt(),
		UpdatedAt: e.UpdatedAt(),
	}

	return r.q.CreateEstablishment(ctx, params)
}

func (r *establishmentRepository) Update(
	ctx context.Context,
	e *domain.Establishment,
) error {
	return nil
}

func (r *establishmentRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}

func (r *establishmentRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Establishment, error) {
	return nil, nil
}

func (r *establishmentRepository) GetBySlug(
	ctx context.Context,
	slug domain.Slug,
) (*domain.Establishment, error) {

	row, err := r.q.GetEstablishmentBySlug(ctx, slug.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEstablishmentNotFound
		}
		return nil, err
	}

	var phones []domain.Phone
	if err := json.Unmarshal(row.Phones, &phones); err != nil {
		return nil, err
	}

	var address domain.Address
	if err := json.Unmarshal(row.Address, &address); err != nil {
		return nil, err
	}

	location, err := domain.NewLocation(row.Lat, row.Lon)
	if err != nil {
		return nil, err
	}

	types := make([]domain.EstablishmentType, len(row.Types))
	for i, t := range row.Types {
		types[i] = domain.EstablishmentType(t)
	}

	est, err := domain.NewEstablishment(
		row.Name,
		row.Slug,
		types,
		location,
		address,
	)
	if err != nil {
		return nil, err
	}

	return est, nil
}

func (r *establishmentRepository) List(
	ctx context.Context,
	filter *domain.EstablishmentFilter,
) ([]*domain.Establishment, error) {
	return nil, nil
}
