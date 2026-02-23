package domain

type EstablishmentFilter struct {
	Types  []string
	Search string
	Limit  int
	Offset int
}
