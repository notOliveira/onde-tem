package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/notOliveira/onde-tem/internal/adapters/dto"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres/sqlc"
	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/ports"
)

type establishmentRepository struct {
	q sqlc.Querier
}

var _ ports.EstablishmentRepository = (*establishmentRepository)(nil)

func NewEstablishmentRepository(q sqlc.Querier) ports.EstablishmentRepository {
	return &establishmentRepository{q: q}
}

func (r *establishmentRepository) Create(
	ctx context.Context,
	e *domain.Establishment,
) error {

	params, err := r.toCreateParams(e)
	if err != nil {
		return err
	}

	id, err := r.q.CreateEstablishment(ctx, params)
	if err != nil {
		return err
	}

	parsedID, err := domain.ParseEstablishmentID(id.String())
	if err != nil {
		return err
	}

	e.SetID(parsedID)

	return nil
}

func (r *establishmentRepository) toCreateParams(
	e *domain.Establishment,
) (sqlc.CreateEstablishmentParams, error) {

	phonesJSON, _ := json.Marshal(dto.PhonesFromDomain(e.Phones()))
	addressJSON, _ := json.Marshal(dto.AddressFromDomain(e.Address()))

	types := make([]string, len(e.EstablishmentTypes()))
	for i, t := range e.EstablishmentTypes() {
		types[i] = t.String()
	}

	return sqlc.CreateEstablishmentParams{
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
	}, nil
}

func (r *establishmentRepository) Update(
	ctx context.Context,
	e *domain.Establishment,
) error {

	parsedUUID, err := uuid.Parse(e.ID().String())
	if err != nil {
		return err
	}

	phonesJSON, err := json.Marshal(dto.PhonesFromDomain(e.Phones()))
	if err != nil {
		return err
	}

	addressJSON, err := json.Marshal(dto.AddressFromDomain(e.Address()))
	if err != nil {
		return err
	}

	types := make([]string, len(e.EstablishmentTypes()))
	for i, t := range e.EstablishmentTypes() {
		types[i] = t.String()
	}

	params := sqlc.UpdateEstablishmentParams{
		ID:        parsedUUID,
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
		UpdatedAt: e.UpdatedAt(),
	}

	if err := r.q.UpdateEstablishment(ctx, params); err != nil {
		return err
	}

	return nil
}

func (r *establishmentRepository) Delete(
	ctx context.Context,
	id domain.EstablishmentID,
) error {

	parsedUUID, err := uuid.Parse(id.String())
	if err != nil {
		return err
	}

	if err := r.q.DeleteEstablishment(ctx, parsedUUID); err != nil {
		return err
	}

	return nil
}

func (r *establishmentRepository) GetByID(
	ctx context.Context,
	id domain.EstablishmentID,
) (*domain.Establishment, error) {

	parsedUUID, err := uuid.Parse(id.String())
	if err != nil {
		return nil, domain.ErrInvalidEstablishmentID
	}

	row, err := r.q.GetEstablishmentByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEstablishmentNotFound
		}
		return nil, err
	}

	dbRow, err := dto.GetEstablishmentByIDRowToDBRow(row)
	if err != nil {
		return nil, err
	}

	return dto.RowToEstablishment(dbRow)
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

	dbRow, err := dto.GetEstablishmentBySlugRowToDBRow(row)
	if err != nil {
		return nil, err
	}

	return dto.RowToEstablishment(dbRow)
}

func (r *establishmentRepository) List(
	ctx context.Context,
	filter *domain.EstablishmentFilter,
) ([]*domain.Establishment, error) {

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := max(filter.Offset, 0)

	params := sqlc.ListEstablishmentsParams{
		Types:  filter.Types,
		Search: filter.Search,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := r.q.ListEstablishments(ctx, params)
	if err != nil {
		return nil, err
	}

	establishments := make([]*domain.Establishment, 0, len(rows))

	for _, row := range rows {
		dbRow, err := dto.ListEstablishmentsRowToDBRow(row)
		if err != nil {
			return nil, err
		}
		est, err := dto.RowToEstablishment(dbRow)
		if err != nil {
			return nil, err
		}
		establishments = append(establishments, est)
	}

	return establishments, nil
}
