package domain

import (
	"strings"
	"time"
)

// Establishment represents a business or location that can be searched for on the platform.
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

// NewEstablishment creates a new Establishment instance with the provided parameters.
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

	newLocation, err := NewLocation(location.lat, location.lon)
	if err != nil {
		return nil, err
	}

	newAddress, err := NewAddress(
		address.street,
		address.number,
		address.district,
		address.city,
		address.state,
		address.country,
		address.zipCode,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Establishment{
		id:        "",
		name:      name,
		slug:      finalSlug,
		types:     types,
		location:  newLocation,
		address:   newAddress,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the unique identifier of the establishment.
func (e *Establishment) ID() EstablishmentID {
	return e.id
}

// Name returns the name of the establishment.
func (e *Establishment) Name() string { return e.name }

// Slug returns the slug of the establishment.
func (e *Establishment) Slug() Slug { return e.slug }

// EstablishmentTypes returns a copy of the slice of establishment types.
func (e *Establishment) EstablishmentTypes() []EstablishmentType {
	return append([]EstablishmentType{}, e.types...)
}

// Email returns the email of the establishment.
func (e *Establishment) Email() string { return e.email }

// Website returns the website of the establishment.
func (e *Establishment) Website() string { return e.website }

// Phones returns a copy of the slice of phones.
func (e *Establishment) Phones() []Phone { return append([]Phone{}, e.phones...) }

// Location returns the location of the establishment.
func (e *Establishment) Location() Location { return e.location }

// Address returns the address of the establishment.
func (e *Establishment) Address() Address { return e.address }

// Timezone returns the timezone of the establishment.
func (e *Establishment) Timezone() string { return e.timezone }

// CreatedAt returns the creation time of the establishment.
func (e *Establishment) CreatedAt() time.Time { return e.createdAt }

// UpdatedAt returns the last update time of the establishment.
func (e *Establishment) UpdatedAt() time.Time { return e.updatedAt }

// UpdateName updates the name of the establishment. The name cannot be empty.
func (e *Establishment) UpdateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrInvalidName
	}
	e.name = name
	e.updatedAt = time.Now()
	return nil
}

// UpdateContact updates the contact information of the establishment. Both email and website cannot be blank.
func (e *Establishment) UpdateContact(email, website string) error {
	email = strings.TrimSpace(email)
	website = strings.TrimSpace(website)

	if email == "" && website == "" {
		return ErrBlankEmailAndWebsite
	}

	e.email = email
	e.website = website
	e.updatedAt = time.Now()
	return nil
}

// AddPhone adds a new phone to the establishment's list of phones.
func (e *Establishment) AddPhone(phone Phone) error {
	newPhone, err := NewPhone(phone.countryCode, phone.number, phone.label)
	if err != nil {
		return err
	}

	e.phones = append(e.phones, newPhone)
	e.updatedAt = time.Now()
	return nil
}

// RemovePhone removes a phone from the establishment's list of phones based on the country code and number.
func (e *Establishment) RemovePhone(countryCode, number string) {
	number = strings.TrimSpace(number)

	for i, p := range e.phones {
		if p.countryCode == countryCode && p.number == number {
			e.phones = append(e.phones[:i], e.phones[i+1:]...)
			e.updatedAt = time.Now()
			return
		}
	}
}

// UpdateLocation updates the location of the establishment.
func (e *Establishment) UpdateLocation(location Location) error {
	if !location.IsValid() {
		return ErrInvalidLocation
	}
	e.location = location
	e.updatedAt = time.Now()
	return nil
}

// UpdateAddress updates the address of the establishment.
func (e *Establishment) UpdateAddress(address Address) error {
	newAddress, err := NewAddress(
		address.street,
		address.number,
		address.district,
		address.city,
		address.state,
		address.country,
		address.zipCode,
	)
	if err != nil {
		return err
	}
	e.address = newAddress
	e.updatedAt = time.Now()
	return nil
}

// UpdateTimezone updates the timezone of the establishment.
func (e *Establishment) UpdateTimezone(timezone string) error {
	_, err := time.LoadLocation(timezone)
	if err != nil {
		return ErrInvalidTimezone
	}
	e.timezone = timezone
	e.updatedAt = time.Now()
	return nil
}

// SetID sets the unique identifier of the establishment. This method is intended to be used by the repository when persisting a new establishment.
func (e *Establishment) SetID(id EstablishmentID) {
	e.id = id
}

// SetTypes updates the types of the establishment. The establishment must have at least one type.
func (e *Establishment) SetTypes(types []EstablishmentType) error {
	if len(types) == 0 {
		return ErrInvalidTypes
	}
	e.types = types
	e.updatedAt = time.Now()
	return nil
}

// RehydrateEstablishment is used strictly by repositories to restore an entity from the database state.
// It skips business validations meant for creation and preserves database timestamps.
func RehydrateEstablishment(
	id EstablishmentID,
	name string,
	slug Slug,
	types []EstablishmentType,
	email string,
	website string,
	phones []Phone,
	location Location,
	address Address,
	timezone string,
	createdAt time.Time,
	updatedAt time.Time,
) *Establishment {
	return &Establishment{
		id:        id,
		name:      name,
		slug:      slug,
		types:     types,
		email:     email,
		website:   website,
		phones:    phones,
		location:  location,
		address:   address,
		timezone:  timezone,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
