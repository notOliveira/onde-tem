package domain

type Location struct {
	Lat float64
	Lon float64
}

func (l Location) IsValid() bool {
	if l.Lat < -90 || l.Lat > 90 {
		return false
	}
	if l.Lon < -180 || l.Lon > 180 {
		return false
	}
	return true
}
