package domain

type EstablishmentType string

const (
	EstablishmentTypeRestaurant EstablishmentType = "restaurant"
)

func (t EstablishmentType) String() string {
	return string(t)
}
