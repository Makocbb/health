package providers

import (
	"errors"
	"math/rand"
	"time"

	"go-take-home-test/internal/models"
)

var ErrLookupFailed = errors.New("ideal postcodes: lookup failed")

type poscodeRepository struct {
}

// LookupPostcode is a mock Ideal Postcodes geocoding API.
// It succeeds ~95% of the time and always sleeps ~1s to simulate latency.
func LookupPostcode(postcode string) (*models.Coordinates, error) {
	time.Sleep(time.Second)

	if rand.Float64() < 0.95 {
		return &models.Coordinates{Longitude: 50.05, Latitude: -5.05}, nil
	}

	return nil, ErrLookupFailed
}
