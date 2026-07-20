package models

import "errors"

var ErrLookupFailed = errors.New("ideal postcodes: lookup failed")

type Coordinates struct {
	Longitude float64
	Latitude  float64
}
