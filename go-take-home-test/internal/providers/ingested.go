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

type ingestedFormRepository struct {
	db *bun.DB
}

func NewIngestedFormRepository(db *bun.DB) repositories.IngestedFormRepository {
	return &ingestedFormRepository{db: db}
}

func (r *ingestedFormRepository) Create(ctx context.Context, item *models.IngestedForm) (*models.IngestedForm, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		slog.Error("failed to create ingested form", "error", err, "form", item)
		return nil, err
	}
	return item, nil
}

func (r *ingestedFormRepository) FindByID(ctx context.Context, id int64) (*models.IngestedForm, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	data := &models.IngestedForm{ID: id}
	err := r.db.NewSelect().Model(data).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrIngestedFormNotFound
		}
		slog.Error("failed to find ingested form by id", "error", err, "id", id)
		return nil, err
	}
	return data, nil
}

func (r *ingestedFormRepository) FindAll(ctx context.Context, params *models.IngestedFormParams) (int, []models.IngestedForm, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	var forms []models.IngestedForm
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
		slog.Error("failed to find all ingested forms", "error", err, "params", params)
		return 0, nil, err
	}
	return count, forms, nil
}

func (r *ingestedFormRepository) Patch(ctx context.Context, item *models.IngestedForm, columns ...string) (*models.IngestedForm, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewUpdate().Model(item).Column(columns...).WherePK().Exec(ctx)
	if err != nil {
		slog.Error("failed to patch ingested form", "error", err, "form", item, "columns", columns)
		return nil, err
	}
	return item, nil
}

func (r *ingestedFormRepository) Update(ctx context.Context, item *models.IngestedForm) (*models.IngestedForm, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewUpdate().Model(item).WherePK().Exec(ctx)
	if err != nil {
		slog.Error("failed to update ingested form", "error", err, "form", item)
		return nil, err
	}
	return item, nil
}

func (r *ingestedFormRepository) Delete(ctx context.Context, id int64) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewDelete().Model((*models.IngestedForm)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		slog.Error("failed to delete ingested form", "error", err, "id", id)
		return err
	}
	return nil
}
