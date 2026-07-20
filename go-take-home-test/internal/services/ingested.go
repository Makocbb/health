package services

import (
	"context"
	"errors"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type ingestedFormService struct {
	repo repositories.IngestedFormRepository
}

type IngestedFormService interface {
	Create(ctx context.Context, form *models.IngestedForm) (*models.IngestedForm, error)
	GetOrCreate(ctx context.Context, form *models.IngestedForm) (created bool, saved *models.IngestedForm, err error)
	GetByID(ctx context.Context, id int64) (*models.IngestedForm, error)
	GetBySessionID(ctx context.Context, sessionID string) (*models.IngestedForm, error)
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
	if form.Status == "" {
		form.Status = models.IngestedStatusPending
	}
	return s.repo.Create(ctx, form)
}

// GetOrCreate returns an existing form with the same content fingerprint,
// or creates a new one when none exists.
func (s *ingestedFormService) GetOrCreate(ctx context.Context, form *models.IngestedForm) (bool, *models.IngestedForm, error) {
	existing, err := s.repo.FindByFingerprint(ctx, form.Fingerprint)
	if err == nil {
		return false, existing, nil
	}
	if !errors.Is(err, models.ErrIngestedFormNotFound) {
		return false, nil, err
	}

	created, err := s.Create(ctx, form)
	if err != nil {
		// Race: another request inserted the same fingerprint.
		existing, findErr := s.repo.FindByFingerprint(ctx, form.Fingerprint)
		if findErr == nil {
			return false, existing, nil
		}
		return false, nil, err
	}
	return true, created, nil
}

func (s *ingestedFormService) GetByID(ctx context.Context, id int64) (*models.IngestedForm, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ingestedFormService) GetBySessionID(ctx context.Context, sessionID string) (*models.IngestedForm, error) {
	return s.repo.FindBySessionID(ctx, sessionID)
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
