package domain

type Phone struct {
	countryCode string
	number      string
	label       string
}

func (p Phone) IsValid() bool {
	return p.number != ""
}
