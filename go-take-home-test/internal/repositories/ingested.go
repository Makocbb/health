package repositories

import (
	"context"

	"go-take-home-test/internal/models"
)

type IngestedFormRepository interface {
	Create(ctx context.Context, data *models.IngestedForm) (*models.IngestedForm, error)
	FindAll(ctx context.Context, params *models.IngestedFormParams) (int, []models.IngestedForm, error)
	FindByID(ctx context.Context, id int64) (*models.IngestedForm, error)
	Patch(ctx context.Context, data *models.IngestedForm, columns ...string) (*models.IngestedForm, error)
	Update(ctx context.Context, data *models.IngestedForm) (*models.IngestedForm, error)
	Delete(ctx context.Context, id int64) error
}
