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

type transformLogRepository struct {
	db *bun.DB
}

func NewTransformLogRepository(db *bun.DB) repositories.TransformLogRepository {
	return &transformLogRepository{db: db}
}

func (r *transformLogRepository) Create(ctx context.Context, item *models.TransformLog) (*models.TransformLog, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		slog.Error("failed to create transform log", "error", err, "log", item)
		return nil, err
	}
	return item, nil
}

func (r *transformLogRepository) FindByID(ctx context.Context, id int64) (*models.TransformLog, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	data := &models.TransformLog{ID: id}
	err := r.db.NewSelect().Model(data).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrTransformLogNotFound
		}
		slog.Error("failed to find transform log by id", "error", err, "id", id)
		return nil, err
	}
	return data, nil
}

func (r *transformLogRepository) FindAll(ctx context.Context, params *models.TransformLogParams) (int, []models.TransformLog, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	var logs []models.TransformLog
	query := r.db.NewSelect().Model(&logs)

	if params.Success != nil {
		query = query.Where("success = ?", *params.Success)
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
		slog.Error("failed to find all transform logs", "error", err, "params", params)
		return 0, nil, err
	}
	return count, logs, nil
}

func (r *transformLogRepository) Patch(ctx context.Context, item *models.TransformLog, columns ...string) (*models.TransformLog, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewUpdate().Model(item).Column(columns...).WherePK().Exec(ctx)
	if err != nil {
		slog.Error("failed to patch transform log", "error", err, "log", item, "columns", columns)
		return nil, err
	}
	return item, nil
}

func (r *transformLogRepository) Update(ctx context.Context, item *models.TransformLog) (*models.TransformLog, error) {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewUpdate().Model(item).WherePK().Exec(ctx)
	if err != nil {
		slog.Error("failed to update transform log", "error", err, "log", item)
		return nil, err
	}
	return item, nil
}

func (r *transformLogRepository) Delete(ctx context.Context, id int64) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	_, err := r.db.NewDelete().Model((*models.TransformLog)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		slog.Error("failed to delete transform log", "error", err, "id", id)
		return err
	}
	return nil
}
