package domain

import "encoding/json"

// Location represents the geographical coordinates (latitude and longitude) of an establishment.
type Location struct {
	lat float64
	lon float64
}

func NewLocation(lat, lon float64) (Location, error) {
	l := Location{lat: lat, lon: lon}
	if !l.IsValid() {
		return Location{}, ErrInvalidLocation
	}
	return l, nil
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

func (l Location) MarshalJSON() ([]byte, error) {
	type Alias Location
	return json.Marshal(struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}{
		Lat: l.lat,
		Lon: l.lon,
	})
}

func (l Location) Lat() float64 { return l.lat }
func (l Location) Lon() float64 { return l.lon }
