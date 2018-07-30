package types

import "github.com/SUN-XIN/know-your-friends/geo"

const (
	PLACE_NAME_HOME   = "home"
	PLACE_NAME_SCHOOL = "school"
	PLACE_NAME_WORK   = "work"
)

type SignificantPlace struct {
	UserID string
	ID     string
	LLBox  geo.LLBox
	Name   string // Home, School/Work, enum ?
}

func (sp *SignificantPlace) IsIn(lat, lng float64) bool {
	return sp.LLBox.Contains(&geo.LatLng{
		Lat: lat,
		Lng: lng,
	})
}
