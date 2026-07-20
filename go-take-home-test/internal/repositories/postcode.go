package repositories

import (
	"context"
	"go-take-home-test/internal/models"
)

type PostcodeRepository interface {
	LookupPostcode(ctx context.Context, postcode string) (*models.Coordinates, error)
}
