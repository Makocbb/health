package app

import (
	"fmt"
	"time"

	"go-take-home-test/internal/controllers"
	"go-take-home-test/internal/models"
	"go-take-home-test/internal/providers"
	"go-take-home-test/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config holds application wiring options.
type Config struct {
	DBPath              string
	MigrationsPath      string
	VersionFilePath     string
	BaseURL             string
	QueueMaxAttempts    int
	QueueInitialBackoff time.Duration
}

var defaultConfig = Config{
	DBPath:              "./data/db.sqlite",
	MigrationsPath:      "./internal/migrations",
	VersionFilePath:     "./migration_version.txt",
	BaseURL:             "http://localhost:8080",
	QueueMaxAttempts:    3,
	QueueInitialBackoff: time.Second,
}

// New builds the Echo application with default config.
func New() *echo.Echo {
	e, err := NewWithConfig(defaultConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create app: %w", err))
	}
	return e
}

// NewWithConfig builds the Echo application with the provided config.
func NewWithConfig(cfg Config) (*echo.Echo, error) {
	if cfg.DBPath == "" {
		cfg.DBPath = defaultConfig.DBPath
	}
	if cfg.MigrationsPath == "" {
		cfg.MigrationsPath = defaultConfig.MigrationsPath
	}
	if cfg.VersionFilePath == "" {
		cfg.VersionFilePath = defaultConfig.VersionFilePath
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultConfig.BaseURL
	}
	if cfg.QueueMaxAttempts <= 0 {
		cfg.QueueMaxAttempts = defaultConfig.QueueMaxAttempts
	}
	if cfg.QueueInitialBackoff <= 0 {
		cfg.QueueInitialBackoff = defaultConfig.QueueInitialBackoff
	}

	migrationOpts := models.NewMigrationOptions(
		models.MigrationWithVersionFilePath(cfg.VersionFilePath),
		models.MigrationWithMigrationsPath(cfg.MigrationsPath),
	)
	queueOpts := models.NewQueueOptions(
		models.QueueWithBaseURL(cfg.BaseURL),
		models.QueueWithMaxAttempts(cfg.QueueMaxAttempts),
		models.QueueWithInitialBackoff(cfg.QueueInitialBackoff),
	)

	dbRepository, err := providers.NewDBRepository(cfg.DBPath, migrationOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create db repository: %w", err)
	}

	db := dbRepository.GetDB()

	ingestedFormRepository := providers.NewIngestedFormRepository(db)
	transformedFormRepository := providers.NewTransformedFormRepository(db)
	transformLogRepository := providers.NewTransformLogRepository(db)
	queueRepository := providers.NewQueueRepository(queueOpts)
	emailRepository := providers.NewEmailRepository()
	postcodeRepository := providers.NewPostcodeRepository()

	ingestedFormService := services.NewIngestedFormService(ingestedFormRepository)
	transformedFormService := services.NewTransformedFormService(transformedFormRepository)
	transformLogService := services.NewTransformLogService(transformLogRepository)
	queueService := services.NewQueueService(queueRepository)
	emailService := services.NewEmailService(emailRepository)
	postcodeService := services.NewPostcodeService(postcodeRepository)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	workerCtl := controllers.NewWorkerController(
		ingestedFormService,
		transformedFormService,
		transformLogService,
		queueService,
		postcodeService,
		emailService,
	)
	workerCtl.Routes(e.Group("/workers"))
	workerCtl.PublicRoutes(e)

	return e, nil
}
