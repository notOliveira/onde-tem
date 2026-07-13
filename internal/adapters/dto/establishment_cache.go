package dto

import (
	"time"

	"github.com/notOliveira/onde-tem/internal/core/domain"
)

type EstablishmentCacheDTO struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	Types     []string    `json:"types"`
	Email     string      `json:"email"`
	Website   string      `json:"website"`
	Phones    []PhoneDTO  `json:"phones"`
	Location  LocationDTO `json:"location"`
	Address   AddressDTO  `json:"address"`
	Timezone  string      `json:"timezone"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

func typeSliceToStrings(types []domain.EstablishmentType) []string {
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = t.String()
	}
	return result
}

func ToEstablishmentCacheDTO(e *domain.Establishment) *EstablishmentCacheDTO {
	return &EstablishmentCacheDTO{
		ID:        e.ID().String(),
		Name:      e.Name(),
		Slug:      e.Slug().String(),
		Types:     typeSliceToStrings(e.EstablishmentTypes()),
		Email:     e.Email(),
		Website:   e.Website(),
		Phones:    PhonesFromDomain(e.Phones()),
		Location:  LocationFromDomain(e.Location()),
		Address:   AddressFromDomain(e.Address()),
		Timezone:  e.Timezone(),
		CreatedAt: e.CreatedAt().Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt().Format(time.RFC3339),
	}
}
