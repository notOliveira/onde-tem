package dto

import (
	"encoding/json"
	"github.com/notOliveira/onde-tem/internal/core/domain"
)

type AddressDTO struct {
	Street   string `json:"street"`
	Number   string `json:"number"`
	District string `json:"district"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	ZipCode  string `json:"zipCode"`
}

func AddressFromDomain(a domain.Address) AddressDTO {
	return AddressDTO{
		Street:   a.Street(),
		Number:   a.Number(),
		District: a.District(),
		City:     a.City(),
		State:    a.State(),
		Country:  a.Country(),
		ZipCode:  a.ZipCode(),
	}
}

func (dto AddressDTO) ToDomain() (domain.Address, error) {
	return domain.NewAddress(dto.Street, dto.Number, dto.District, dto.City, dto.State, dto.Country, dto.ZipCode)
}

func MarshalAddress(a domain.Address) ([]byte, error) {
	return json.Marshal(AddressFromDomain(a))
}

func UnmarshalAddress(data []byte) (domain.Address, error) {
	var dto AddressDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return domain.Address{}, err
	}
	return dto.ToDomain()
}
