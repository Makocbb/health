package providers

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"

	"github.com/uptrace/bun"
)

type transformedFormRepository struct {
	db *bun.DB
}

func NewTransformedFormRepository(db *bun.DB) repositories.TransformedFormRepository {
	return &transformedFormRepository{db: db}
}

func (r *transformedFormRepository) Create(ctx context.Context, item *models.TransformedForm) (*models.TransformedForm, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		slog.Error("failed to create transformed form", "error", err, "form", item)
		return nil, err
	}
	return item, nil
}

func (r *transformedFormRepository) FindByID(ctx context.Context, id int64) (*models.TransformedForm, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	data := &models.TransformedForm{ID: id}
	err := r.db.NewSelect().Model(data).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrTransformedFormNotFound
		}
		slog.Error("failed to find transformed form by id", "error", err, "id", id)
		return nil, err
	}
	return data, nil
}

func (r *transformedFormRepository) FindAll(ctx context.Context, params *models.TransformedFormParams) (int, []models.TransformedForm, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	var forms []models.TransformedForm
	query := r.db.NewSelect().Model(&forms)

	if params.SessionID != "" {
		query = query.Where("session_id = ?", params.SessionID)
	}
	if params.ApplicationReference != "" {
		query = query.Where("application_reference = ?", params.ApplicationReference)
	}

	perPage := params.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	page := params.Page
	if page < 1 {
		page = 1
	}

	count, err := query.Limit(perPage).Offset(perPage * (page - 1)).ScanAndCount(ctx)
	if err != nil {
		slog.Error("failed to find all transformed forms", "error", err, "params", params)
		return 0, nil, err
	}
	return count, forms, nil
}

func (r *transformedFormRepository) Patch(ctx context.Context, item *models.TransformedForm, columns ...string) (*models.TransformedForm, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewUpdate().Model(item).Column(columns...).WherePK().Exec(ctx)
	if err != nil {
		slog.Error("failed to patch transformed form", "error", err, "form", item, "columns", columns)
		return nil, err
	}
	return item, nil
}

func (r *transformedFormRepository) Update(ctx context.Context, item *models.TransformedForm) (*models.TransformedForm, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewUpdate().Model(item).WherePK().Exec(ctx)
	if err != nil {
		slog.Error("failed to update transformed form", "error", err, "form", item)
		return nil, err
	}
	return item, nil
}

func (r *transformedFormRepository) Delete(ctx context.Context, id int64) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewDelete().Model((*models.TransformedForm)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		slog.Error("failed to delete transformed form", "error", err, "id", id)
		return err
	}
	return nil
}
