package domain

import "errors"

var (
	ErrInvalidName           = errors.New("invalid establishment name")
	ErrInvalidTypes          = errors.New("establishment must have at least one type")
	ErrInvalidLocation       = errors.New("invalid location")
	ErrInvalidPhone          = errors.New("invalid phone number")
	ErrInvalidAddress        = errors.New("invalid address")
	ErrInvalidTimezone       = errors.New("invalid timezone")
	ErrBlankEmailAndWebsite  = errors.New("both email and website cannot be blank")
	ErrEstablishmentNotFound = errors.New("establishment not found")
)
