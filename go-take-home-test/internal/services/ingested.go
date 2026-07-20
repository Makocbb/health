package services

import (
	"context"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type ingestedFormService struct {
	repo repositories.IngestedFormRepository
}

type IngestedFormService interface {
	Create(ctx context.Context, form *models.IngestedForm) (*models.IngestedForm, error)
	GetByID(ctx context.Context, id int64) (*models.IngestedForm, error)
	GetAll(ctx context.Context, params *models.IngestedFormParams) (int, []models.IngestedForm, error)
	Patch(ctx context.Context, form *models.IngestedForm, columns ...string) (*models.IngestedForm, error)
	Update(ctx context.Context, form *models.IngestedForm) (*models.IngestedForm, error)
	Delete(ctx context.Context, id int64) error
}

func NewIngestedFormService(repo repositories.IngestedFormRepository) IngestedFormService {
	return &ingestedFormService{repo: repo}
}

func (s *ingestedFormService) Create(ctx context.Context, form *models.IngestedForm) (*models.IngestedForm, error) {
	now := time.Now()
	form.CreatedAt = now
	form.UpdatedAt = now
	return s.repo.Create(ctx, form)
}

func (s *ingestedFormService) GetByID(ctx context.Context, id int64) (*models.IngestedForm, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ingestedFormService) GetAll(ctx context.Context, params *models.IngestedFormParams) (int, []models.IngestedForm, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *ingestedFormService) Patch(ctx context.Context, form *models.IngestedForm, columns ...string) (*models.IngestedForm, error) {
	form.UpdatedAt = time.Now()
	columns = append(columns, "updated_at")
	return s.repo.Patch(ctx, form, columns...)
}

func (s *ingestedFormService) Update(ctx context.Context, form *models.IngestedForm) (*models.IngestedForm, error) {
	form.UpdatedAt = time.Now()
	return s.repo.Update(ctx, form)
}

func (s *ingestedFormService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
