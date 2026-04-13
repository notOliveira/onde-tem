package dto

import "github.com/notOliveira/onde-tem/internal/core/domain"

type LocationDTO struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func LocationFromDomain(l domain.Location) LocationDTO {
	return LocationDTO{
		Lat: l.Lat(),
		Lon: l.Lon(),
	}
}

func (dto LocationDTO) ToDomain() (domain.Location, error) {
	return domain.NewLocation(dto.Lat, dto.Lon)
}
