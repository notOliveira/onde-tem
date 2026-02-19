package domain

type Phone struct {
	CountryCode string
	Number      string
	Label       string
}

func (p Phone) IsValid() bool {
	return p.Number != ""
}
