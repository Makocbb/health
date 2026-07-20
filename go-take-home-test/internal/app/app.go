package app

import (
	"go-take-home-test/internal/models"
	"go-take-home-test/internal/providers"
	"go-take-home-test/internal/services"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// appConfig is the configuration for the application.
// Normally I would put it elsewhere, but given the scope of the task,
// I will just keep it here. Its going to be convenient to do quick edits.
type appConfig struct {
	dbPath          string
	migrationsPath  string
	versionFilePath string
	baseURL         string
}

var defaultAppConfig = appConfig{
	dbPath:          "./data/db.sqlite",
	migrationsPath:  "./internal/migrations",
	versionFilePath: "./migration_version.txt",
	baseURL:         "http://localhost:8080",
}

// New builds the Echo application with the starter /ingest route.
func New() *echo.Echo {
	// Options
	migrationOpts := models.NewMigrationOptions(
		models.MigrationWithVersionFilePath(defaultAppConfig.versionFilePath),
		models.MigrationWithMigrationsPath(defaultAppConfig.migrationsPath),
	)
	queueOpts := models.NewQueueOptions(
		models.QueueWithBaseURL(defaultAppConfig.baseURL),
	)

	// Repositories
	dbRepository, err := providers.NewDBRepository(defaultAppConfig.dbPath, migrationOpts)
	if err != nil {
		log.Fatalf("failed to create db repository: %v", err)
		return nil
	}

	db := dbRepository.GetDB()

	var (
		// database repositories
		ingestedFormRepository    = providers.NewIngestedFormRepository(db)
		transformedFormRepository = providers.NewTransformedFormRepository(db)
		transformLogRepository    = providers.NewTransformLogRepository(db)

		// queue repository
		queueRepository = providers.NewQueueRepository(queueOpts)
	)

	// Services
	var (
		// data services
		ingestedFormService    = services.NewIngestedFormService(ingestedFormRepository)
		transformedFormService = services.NewTransformedFormService(transformedFormRepository)
		transformLogService    = services.NewTransformLogService(transformLogRepository)

		// queue service
		queueService = services.NewQueueService(queueRepository)
	)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.POST("/ingest", ingest)
	return e
}

func ingest(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Ingesting form data",
	})
}
