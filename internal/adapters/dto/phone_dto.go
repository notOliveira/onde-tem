package dto

import "github.com/notOliveira/onde-tem/internal/core/domain"

type PhoneDTO struct {
	CountryCode string `json:"countryCode"`
	Number      string `json:"number"`
	Label       string `json:"label"`
}

func PhoneFromDomain(p domain.Phone) PhoneDTO {
	return PhoneDTO{
		CountryCode: p.CountryCode(),
		Number:      p.Number(),
		Label:       p.Label(),
	}
}

func PhonesFromDomain(phones []domain.Phone) []PhoneDTO {
	dtos := make([]PhoneDTO, len(phones))
	for i, p := range phones {
		dtos[i] = PhoneFromDomain(p)
	}
	return dtos
}
