package providers

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/repositories"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

// writeMu is a package-level mutex shared across all repositories to serialise
// all SQLite write operations and prevent SQLITE_BUSY (database is locked) errors.
var writeMu sync.RWMutex

type dbRepository struct {
	db   *bun.DB
	opts models.MigrationOptions
}

// NewDBRepository opens the SQLite database and applies pending SQL migrations.
func NewDBRepository(path string, opts ...models.MigrationOption) (repositories.DBRepository, error) {
	options := models.MigrationOptions{
		VersionFilePath: "./migration_version.txt",
		MigrationsPath:  "./internal/migrations",
	}
	for _, opt := range opts {
		opt(&options)
	}

	repo := &dbRepository{opts: options}

	db, err := repo.newBunDb(path)
	if err != nil {
		slog.Error("failed to create bun database", "error", err)
		return nil, err
	}
	repo.db = db

	if err := repo.applyMigrations(context.Background()); err != nil {
		_ = db.Close()
		slog.Error("failed to apply migrations", "error", err)
		return nil, err
	}

	return repo, nil
}

func (r *dbRepository) GetDB() *bun.DB {
	return r.db
}

func (r *dbRepository) newBunDb(path string) (*bun.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	// Only create the database file if it doesn't already exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	dsn := fmt.Sprintf(
		"file:%s?mode=rw&_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=30000",
		path,
	)

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, err
	}

	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}

func (r *dbRepository) GetCurrentVersion(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "0", nil
		}
		return "", fmt.Errorf("failed to read version file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

func (r *dbRepository) UpdateVersion(filePath string, version string) error {
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(filePath, []byte(version), 0644); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}

	return nil
}

// applyMigrations reads the current version from file and applies all newer migrations in sequence.
// If the version file doesn't exist, assumes a fresh install (version 0) and applies all migrations.
func (r *dbRepository) applyMigrations(ctx context.Context) error {
	currentVersion, err := r.GetCurrentVersion(r.opts.VersionFilePath)
	if err != nil {
		slog.Debug("Failed to read current version, starting from 0", "error", err)
		currentVersion = "0"
	}

	slog.Info("Current migration version", "version", currentVersion)

	entries, err := os.ReadDir(r.opts.MigrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrations := extractAndSortMigrations(entries, currentVersion)

	if len(migrations) == 0 {
		versionFileExists, err := r.versionFileExists()
		if err == nil && !versionFileExists && currentVersion == "0" {
			latestVersion := getLatestMigrationVersion(entries)
			if latestVersion != "" {
				if err := r.UpdateVersion(r.opts.VersionFilePath, latestVersion); err != nil {
					slog.Warn("Failed to update version file after applying all migrations", "error", err)
				} else {
					slog.Info("Fresh install: set migration version file to latest", "version", latestVersion)
				}
			} else if err := r.UpdateVersion(r.opts.VersionFilePath, "0"); err != nil {
				slog.Warn("Failed to create version file for fresh install", "error", err)
			} else {
				slog.Info("Fresh install: no migrations found, created version file", "version", "0")
			}
		}
		slog.Info("Migrations up to date", "currentVersion", currentVersion, "reason", "no pending migrations")
		return nil
	}

	slog.Info("Found pending migrations", "count", len(migrations), "migrations", strings.Join(migrations, ", "))

	appliedMigrations := make([]string, 0, len(migrations))
	for _, migration := range migrations {
		if err := r.applyMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration, err)
		}
		appliedMigrations = append(appliedMigrations, migration)
	}

	if len(appliedMigrations) > 0 {
		latestVersion := extractVersionFromFilename(appliedMigrations[len(appliedMigrations)-1])
		slog.Info(
			"All pending migrations applied successfully",
			"count", len(appliedMigrations),
			"migrations", strings.Join(appliedMigrations, ", "),
			"newVersion", latestVersion,
		)
	}

	return nil
}

func (r *dbRepository) applyMigration(ctx context.Context, migration string) error {
	version := extractVersionFromFilename(migration)
	if version == "" {
		return fmt.Errorf("invalid migration filename: %s", migration)
	}

	migrationPath := filepath.Join(r.opts.MigrationsPath, migration)

	sqlContent, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	slog.Info("Starting migration application", "file", migration, "version", version)

	if _, err = r.db.ExecContext(ctx, string(sqlContent)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	if err := r.UpdateVersion(r.opts.VersionFilePath, version); err != nil {
		return fmt.Errorf("failed to update version file: %w", err)
	}

	slog.Info("Migration completed", "file", migration, "version", version)
	return nil
}

func extractAndSortMigrations(entries []fs.DirEntry, currentVersion string) []string {
	var migrations []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		version := extractVersionFromFilename(name)
		if version == "" {
			continue
		}

		if compareVersions(version, currentVersion) > 0 {
			migrations = append(migrations, name)
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		versionI := extractVersionFromFilename(migrations[i])
		versionJ := extractVersionFromFilename(migrations[j])
		return compareVersions(versionI, versionJ) < 0
	})

	return migrations
}

// extractVersionFromFilename extracts the version timestamp from a migration filename
// e.g., "20260310151606_init.up.sql" -> "20260310151606"
func extractVersionFromFilename(filename string) string {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// compareVersions compares two version strings numerically.
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
func compareVersions(v1, v2 string) int {
	num1, err1 := strconv.ParseInt(v1, 10, 64)
	num2, err2 := strconv.ParseInt(v2, 10, 64)

	if err1 != nil || err2 != nil {
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	}

	if num1 < num2 {
		return -1
	} else if num1 > num2 {
		return 1
	}
	return 0
}

func (r *dbRepository) versionFileExists() (bool, error) {
	_, err := os.Stat(r.opts.VersionFilePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getLatestMigrationVersion(entries []fs.DirEntry) string {
	var latestVersion string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		version := extractVersionFromFilename(name)
		if version == "" {
			continue
		}

		if latestVersion == "" || compareVersions(version, latestVersion) > 0 {
			latestVersion = version
		}
	}

	return latestVersion
}
