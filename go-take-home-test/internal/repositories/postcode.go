package repositories

import "go-take-home-test/internal/models"

type PostcodeRepository interface {
	LookupPostcode(postcode string) (*models.Coordinates, error)
}
