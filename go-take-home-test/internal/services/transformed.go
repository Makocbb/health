package services

import (
	"context"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type transformedFormService struct {
	repo repositories.TransformedFormRepository
}

type TransformedFormService interface {
	Create(ctx context.Context, form *models.TransformedForm) (*models.TransformedForm, error)
	GetByID(ctx context.Context, id int64) (*models.TransformedForm, error)
	GetByIngestedFormID(ctx context.Context, ingestedFormID int64) (*models.TransformedForm, error)
	GetAll(ctx context.Context, params *models.TransformedFormParams) (int, []models.TransformedForm, error)
	Patch(ctx context.Context, form *models.TransformedForm, columns ...string) (*models.TransformedForm, error)
	Update(ctx context.Context, form *models.TransformedForm) (*models.TransformedForm, error)
	Delete(ctx context.Context, id int64) error
}

func NewTransformedFormService(repo repositories.TransformedFormRepository) TransformedFormService {
	return &transformedFormService{repo: repo}
}

func (s *transformedFormService) Create(ctx context.Context, form *models.TransformedForm) (*models.TransformedForm, error) {
	now := time.Now()
	form.CreatedAt = now
	form.UpdatedAt = now
	return s.repo.Create(ctx, form)
}

func (s *transformedFormService) GetByID(ctx context.Context, id int64) (*models.TransformedForm, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *transformedFormService) GetByIngestedFormID(ctx context.Context, ingestedFormID int64) (*models.TransformedForm, error) {
	return s.repo.FindByIngestedFormID(ctx, ingestedFormID)
}

func (s *transformedFormService) GetAll(ctx context.Context, params *models.TransformedFormParams) (int, []models.TransformedForm, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *transformedFormService) Patch(ctx context.Context, form *models.TransformedForm, columns ...string) (*models.TransformedForm, error) {
	form.UpdatedAt = time.Now()
	columns = append(columns, "updated_at")
	return s.repo.Patch(ctx, form, columns...)
}

func (s *transformedFormService) Update(ctx context.Context, form *models.TransformedForm) (*models.TransformedForm, error) {
	form.UpdatedAt = time.Now()
	return s.repo.Update(ctx, form)
}

func (s *transformedFormService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
