package domain

import "errors"

var ErrInvalidEstablishmentID = errors.New("invalid establishment id")

type EstablishmentID string

func NewEstablishmentID(value string) (EstablishmentID, error) {
	if value == "" {
		return "", ErrInvalidEstablishmentID
	}
	return EstablishmentID(value), nil
}

func (id EstablishmentID) String() string {
	return string(id)
}
