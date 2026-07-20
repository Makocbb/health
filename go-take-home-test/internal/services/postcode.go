package services

import (
	"context"
	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type postcodeService struct {
	postcodeRepository repositories.PostcodeRepository
}

type PostcodeService interface {
	LookupPostcode(ctx context.Context, postcode string) (*models.Coordinates, error)
}

func NewPostcodeService(postcodeRepository repositories.PostcodeRepository) PostcodeService {
	return &postcodeService{postcodeRepository: postcodeRepository}
}

func (s *postcodeService) LookupPostcode(ctx context.Context, postcode string) (*models.Coordinates, error) {
	return s.postcodeRepository.LookupPostcode(ctx, postcode)
}
