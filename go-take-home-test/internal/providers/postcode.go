package providers

import (
	"context"
	"math/rand"
	"time"

	"go-take-home-test/internal/models"
)

type poscodeRepository struct {
}

// LookupPostcode is a mock Ideal Postcodes geocoding API.
// It succeeds ~95% of the time and always sleeps ~1s to simulate latency.
func LookupPostcode(ctx context.Context, postcode string) (*models.Coordinates, error) {
	time.Sleep(time.Second)

	if rand.Float64() < 0.95 {
		return &models.Coordinates{Longitude: 50.05, Latitude: -5.05}, nil
	}

	return nil, models.ErrLookupFailed
}
