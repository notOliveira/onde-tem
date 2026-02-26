package domain

import "errors"

var (
	ErrInvalidName           = errors.New("invalid establishment name")
	ErrInvalidTypes          = errors.New("establishment must have at least one type")
	ErrInvalidLocation       = errors.New("invalid location")
	ErrInvalidAddress        = errors.New("invalid address")
	ErrEstablishmentNotFound = errors.New("establishment not found")
)
