package domain

import (
	"encoding/json"
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

func (a Address) MarshalJSON() ([]byte, error) {
	type Alias Address
	return json.Marshal(struct {
		Street   string `json:"street"`
		Number   string `json:"number"`
		District string `json:"district"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		ZipCode  string `json:"zipCode"`
	}{
		Street:   a.street,
		Number:   a.number,
		District: a.district,
		City:     a.city,
		State:    a.state,
		Country:  a.country,
		ZipCode:  a.zipCode,
	})
}

func (a Address) Street() string   { return a.street }
func (a Address) Number() string   { return a.number }
func (a Address) District() string { return a.district }
func (a Address) City() string     { return a.city }
func (a Address) State() string    { return a.state }
func (a Address) Country() string  { return a.country }
func (a Address) ZipCode() string  { return a.zipCode }
