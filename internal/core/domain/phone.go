package domain

import (
	"strings"
)

// Phone represents a phone number associated with an establishment, including an optional country code and label.
type Phone struct {
	countryCode string
	number      string
	label       string
}

func (p Phone) isValid() bool {
	return p.number != ""
}

func NewPhone(countryCode, number, label string) (Phone, error) {
	countryCode = strings.TrimSpace(countryCode)
	number = strings.TrimSpace(number)
	label = strings.TrimSpace(label)

	phone := Phone{
		countryCode: countryCode,
		number:      number,
		label:       label,
	}

	if !phone.isValid() {
		return Phone{}, ErrInvalidPhone
	}

	return phone, nil
}

func (p Phone) CountryCode() string { return p.countryCode }
func (p Phone) Number() string      { return p.number }
func (p Phone) Label() string       { return p.label }
