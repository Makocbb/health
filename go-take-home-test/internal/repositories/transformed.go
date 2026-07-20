package repositories

import (
	"context"

	"go-take-home-test/internal/models"
)

type TransformedFormRepository interface {
	Create(ctx context.Context, data *models.TransformedForm) (*models.TransformedForm, error)
	FindAll(ctx context.Context, params *models.TransformedFormParams) (int, []models.TransformedForm, error)
	FindByID(ctx context.Context, id int64) (*models.TransformedForm, error)
	FindBySessionID(ctx context.Context, sessionID string) (*models.TransformedForm, error)
	FindByIngestedFormID(ctx context.Context, ingestedFormID int64) (*models.TransformedForm, error)
	Patch(ctx context.Context, data *models.TransformedForm, columns ...string) (*models.TransformedForm, error)
	Update(ctx context.Context, data *models.TransformedForm) (*models.TransformedForm, error)
	Delete(ctx context.Context, id int64) error
}
