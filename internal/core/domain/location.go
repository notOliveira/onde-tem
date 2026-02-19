package domain

type Location struct {
	lat float64
	lon float64
}

func (l Location) IsValid() bool {
	if l.lat < -90 || l.lat > 90 {
		return false
	}
	if l.lon < -180 || l.lon > 180 {
		return false
	}
	return true
}
