package domain

import (
	"encoding/json"
	"errors"
)

var ErrInvalidEstablishmentID = errors.New("invalid establishment id")

type EstablishmentID string

func ParseEstablishmentID(value string) (EstablishmentID, error) {
	if value == "" {
		return "", ErrInvalidEstablishmentID
	}
	return EstablishmentID(value), nil
}

func (id EstablishmentID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id EstablishmentID) String() string {
	return string(id)
}
