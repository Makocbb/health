package providers

import (
	"context"
	"math/rand"
	"os"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type postcodeRepository struct {
}

func NewPostcodeRepository() repositories.PostcodeRepository {
	return &postcodeRepository{}
}

// LookupPostcode is a mock Ideal Postcodes geocoding API.
// HEALTH_MOCK_RELIABLE:
//   - "1": always succeed, no sleep
//   - "0": always fail, no sleep
//   - unset/other: ~95% success with ~1s latency
func (r *postcodeRepository) LookupPostcode(ctx context.Context, postcode string) (*models.Coordinates, error) {
	switch os.Getenv("HEALTH_MOCK_RELIABLE") {
	case "1":
		return &models.Coordinates{Longitude: 50.05, Latitude: -5.05}, nil
	case "0":
		return nil, models.ErrLookupFailed
	}

	time.Sleep(time.Second)

	if rand.Float64() < 0.95 {
		return &models.Coordinates{Longitude: 50.05, Latitude: -5.05}, nil
	}

	return nil, models.ErrLookupFailed
}
