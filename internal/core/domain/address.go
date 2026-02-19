package domain

type Address struct {
	street   string
	number   string
	district string
	city     string
	state    string
	country  string
	zipCode  string
}

func (a Address) IsValid() bool {
	return a.street != "" &&
		a.city != "" &&
		a.state != "" &&
		a.country != ""
}
