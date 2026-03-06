package domain

import (
	"strings"
)

// Address represents the physical address of an establishment, including street, number, district, city, state, country, and zip code.
type Address struct {
	street   string
	number   string
	district string
	city     string
	state    string
	country  string
	zipCode  string
}

func (a Address) isValid() bool {
	return a.street != "" &&
		a.city != "" &&
		a.state != "" &&
		a.country != "" &&
		a.zipCode != ""
}

func NewAddress(street, number, district, city, state, country, zipCode string) (Address, error) {
	street = strings.TrimSpace(street)
	number = strings.TrimSpace(number)
	district = strings.TrimSpace(district)
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)
	country = strings.TrimSpace(country)
	zipCode = strings.TrimSpace(zipCode)

	address := Address{
		street:   street,
		number:   number,
		district: district,
		city:     city,
		state:    state,
		country:  country,
		zipCode:  zipCode,
	}

	if !address.isValid() {
		return Address{}, ErrInvalidAddress
	}

	return address, nil
}
