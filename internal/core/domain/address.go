package domain

type Address struct {
	Street   string
	Number   string
	District string
	City     string
	State    string
	Country  string
	ZipCode  string
}

func (a Address) IsValid() bool {
	return a.Street != "" &&
		a.City != "" &&
		a.State != "" &&
		a.Country != ""
}
