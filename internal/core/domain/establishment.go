package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Establishment struct {
	ID       uuid.UUID
	Name     string
	Slug     Slug
	Types    []EstablishmentType
	Email    string
	Website  string
	Timezone string

	Phones   []Phone
	Location Location
	Address  Address

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEstablishment(
	name string,
	rawSlug string,
	types []EstablishmentType,
	location Location,
	address Address,
) (*Establishment, error) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrInvalidName
	}

	var finalSlug Slug
	var err error

	if strings.TrimSpace(rawSlug) == "" {
		finalSlug = GenerateSlugFromName(name)
	} else {
		finalSlug, err = NewSlug(rawSlug)
		if err != nil {
			return nil, err
		}
	}

	if len(types) == 0 {
		return nil, ErrInvalidTypes
	}

	if !location.IsValid() {
		return nil, ErrInvalidLocation
	}

	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}

	now := time.Now()

	return &Establishment{
		ID:        uuid.New(),
		Name:      name,
		Slug:      finalSlug,
		Types:     types,
		Location:  location,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
