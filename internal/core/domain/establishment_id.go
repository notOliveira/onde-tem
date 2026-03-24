package domain

import "errors"

var ErrInvalidEstablishmentID = errors.New("invalid establishment id")

type EstablishmentID string

// TODO: Remover essa função
func NewEstablishmentID(value string) (EstablishmentID, error) {
	if value == "" {
		return "", ErrInvalidEstablishmentID
	}
	return EstablishmentID(value), nil
}

func (id EstablishmentID) String() string {
	return string(id)
}
