package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/notOliveira/onde-tem/internal/adapters/outbound/persistence/postgres/sqlc"
	"github.com/notOliveira/onde-tem/internal/core/domain"
)

type EstablishmentDBRow struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	Types     []string
	Email     string
	Website   string
	Phones    []PhoneDTO
	Lat       float64
	Lon       float64
	Street    string
	Number    string
	District  string
	City      string
	State     string
	Country   string
	ZipCode   string
	Timezone  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func RowToEstablishment(d *EstablishmentDBRow) (*domain.Establishment, error) {
	slug, err := domain.NewSlug(d.Slug)
	if err != nil {
		return nil, err
	}

	types := make([]domain.EstablishmentType, len(d.Types))
	for i, t := range d.Types {
		types[i] = domain.EstablishmentType(t)
	}

	location, err := domain.NewLocation(d.Lat, d.Lon)
	if err != nil {
		return nil, err
	}

	address, err := domain.NewAddress(
		d.Street,
		d.Number,
		d.District,
		d.City,
		d.State,
		d.Country,
		d.ZipCode,
	)
	if err != nil {
		return nil, err
	}

	phones := make([]domain.Phone, 0, len(d.Phones))
	for _, ph := range d.Phones {
		phone, err := domain.NewPhone(ph.CountryCode, ph.Number, ph.Label)
		if err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	id, err := domain.ParseEstablishmentID(d.ID.String())
	if err != nil {
		return nil, err
	}

	return domain.RehydrateEstablishment(
		id,
		d.Name,
		slug,
		types,
		d.Email,
		d.Website,
		phones,
		location,
		address,
		d.Timezone,
		d.CreatedAt,
		d.UpdatedAt,
	), nil
}

func GetEstablishmentByIDRowToDBRow(r sqlc.GetEstablishmentByIDRow) (*EstablishmentDBRow, error) {
	phones, err := unmarshalPhones(r.Phones)
	if err != nil {
		return nil, err
	}
	address, err := unmarshalAddress(r.Address)
	if err != nil {
		return nil, err
	}
	return &EstablishmentDBRow{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		Types:     r.Types,
		Email:     r.Email,
		Website:   r.Website,
		Phones:    phones,
		Lat:       r.Lat,
		Lon:       r.Lon,
		Street:    address.Street,
		Number:    address.Number,
		District:  address.District,
		City:      address.City,
		State:     address.State,
		Country:   address.Country,
		ZipCode:   address.ZipCode,
		Timezone:  r.Timezone,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func GetEstablishmentBySlugRowToDBRow(r sqlc.GetEstablishmentBySlugRow) (*EstablishmentDBRow, error) {
	phones, err := unmarshalPhones(r.Phones)
	if err != nil {
		return nil, err
	}
	address, err := unmarshalAddress(r.Address)
	if err != nil {
		return nil, err
	}
	return &EstablishmentDBRow{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		Types:     r.Types,
		Email:     r.Email,
		Website:   r.Website,
		Phones:    phones,
		Lat:       r.Lat,
		Lon:       r.Lon,
		Street:    address.Street,
		Number:    address.Number,
		District:  address.District,
		City:      address.City,
		State:     address.State,
		Country:   address.Country,
		ZipCode:   address.ZipCode,
		Timezone:  r.Timezone,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func ListEstablishmentsRowToDBRow(r sqlc.ListEstablishmentsRow) (*EstablishmentDBRow, error) {
	phones, err := unmarshalPhones(r.Phones)
	if err != nil {
		return nil, err
	}
	address, err := unmarshalAddress(r.Address)
	if err != nil {
		return nil, err
	}
	return &EstablishmentDBRow{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		Types:     r.Types,
		Email:     r.Email,
		Website:   r.Website,
		Phones:    phones,
		Lat:       r.Lat,
		Lon:       r.Lon,
		Street:    address.Street,
		Number:    address.Number,
		District:  address.District,
		City:      address.City,
		State:     address.State,
		Country:   address.Country,
		ZipCode:   address.ZipCode,
		Timezone:  r.Timezone,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func unmarshalPhones(data []byte) ([]PhoneDTO, error) {
	if len(data) == 0 {
		return []PhoneDTO{}, nil
	}
	var phones []PhoneDTO
	if err := json.Unmarshal(data, &phones); err != nil {
		return nil, err
	}
	return phones, nil
}

func unmarshalAddress(data []byte) (AddressDTO, error) {
	if len(data) == 0 {
		return AddressDTO{}, nil
	}
	var address AddressDTO
	if err := json.Unmarshal(data, &address); err != nil {
		return AddressDTO{}, err
	}
	return address, nil
}
