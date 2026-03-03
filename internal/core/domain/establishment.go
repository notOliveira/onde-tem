package domain

import (
	"strings"
	"time"
)

type Establishment struct {
	id        EstablishmentID
	name      string
	slug      Slug
	types     []EstablishmentType
	email     string
	website   string
	phones    []Phone
	location  Location
	address   Address
	timezone  string
	createdAt time.Time
	updatedAt time.Time
}

func NewEstablishment(
	id EstablishmentID,
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
		id:        id,
		name:      name,
		slug:      finalSlug,
		types:     types,
		location:  location,
		address:   address,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (e *Establishment) ID() EstablishmentID {
	return e.id
}

func (e *Establishment) Name() string {
	return e.name
}

func (e *Establishment) Slug() Slug {
	return e.slug
}

func (e *Establishment) EstablishmentTypes() []EstablishmentType {
	return append([]EstablishmentType{}, e.types...)
}

func (e *Establishment) Email() string {
	return e.email
}

func (e *Establishment) Website() string {
	return e.website
}

func (e *Establishment) Phones() []Phone {
	return append([]Phone{}, e.phones...)
}

func (e *Establishment) Location() Location {
	return e.location
}

func (e *Establishment) Address() Address {
	return e.address
}

func (e *Establishment) Timezone() string {
	return e.timezone
}

func (e *Establishment) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Establishment) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Establishment) UpdateContact(email, website string) {
	e.email = strings.TrimSpace(email)
	e.website = strings.TrimSpace(website)
	e.updatedAt = time.Now()
}

func (e *Establishment) AddPhone(phone Phone) {
	e.phones = append(e.phones, phone)
	e.updatedAt = time.Now()
}

func (e *Establishment) UpdateTimezone(timezone string) {
	e.timezone = timezone
	e.updatedAt = time.Now()
}

func (e *Establishment) Types() []EstablishmentType {
	return append([]EstablishmentType{}, e.types...)
}

func (e *Establishment) SetTypes(types []EstablishmentType) {
	e.types = types
	e.updatedAt = time.Now()
}
