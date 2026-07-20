package repositories

import (
	"context"

	"go-take-home-test/internal/models"
)

type TransformLogRepository interface {
	Create(ctx context.Context, data *models.TransformLog) (*models.TransformLog, error)
	FindAll(ctx context.Context, params *models.TransformLogParams) (int, []models.TransformLog, error)
	FindByID(ctx context.Context, id int64) (*models.TransformLog, error)
	Patch(ctx context.Context, data *models.TransformLog, columns ...string) (*models.TransformLog, error)
	Update(ctx context.Context, data *models.TransformLog) (*models.TransformLog, error)
	Delete(ctx context.Context, id int64) error
}
