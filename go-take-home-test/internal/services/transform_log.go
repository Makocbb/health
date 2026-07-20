package services

import (
	"context"
	"time"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"
)

type transformLogService struct {
	repo repositories.TransformLogRepository
}

type TransformLogService interface {
	Create(ctx context.Context, item *models.TransformLog) (*models.TransformLog, error)
	GetByID(ctx context.Context, id int64) (*models.TransformLog, error)
	GetAll(ctx context.Context, params *models.TransformLogParams) (int, []models.TransformLog, error)
	Patch(ctx context.Context, item *models.TransformLog, columns ...string) (*models.TransformLog, error)
	Update(ctx context.Context, item *models.TransformLog) (*models.TransformLog, error)
	Delete(ctx context.Context, id int64) error
}

func NewTransformLogService(repo repositories.TransformLogRepository) TransformLogService {
	return &transformLogService{repo: repo}
}

func (s *transformLogService) Create(ctx context.Context, item *models.TransformLog) (*models.TransformLog, error) {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	return s.repo.Create(ctx, item)
}

func (s *transformLogService) GetByID(ctx context.Context, id int64) (*models.TransformLog, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *transformLogService) GetAll(ctx context.Context, params *models.TransformLogParams) (int, []models.TransformLog, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *transformLogService) Patch(ctx context.Context, item *models.TransformLog, columns ...string) (*models.TransformLog, error) {
	item.UpdatedAt = time.Now()
	columns = append(columns, "updated_at")
	return s.repo.Patch(ctx, item, columns...)
}

func (s *transformLogService) Update(ctx context.Context, item *models.TransformLog) (*models.TransformLog, error) {
	item.UpdatedAt = time.Now()
	return s.repo.Update(ctx, item)
}

func (s *transformLogService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
