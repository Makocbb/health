package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// New builds the Echo application with the starter /ingest route.
func New() *echo.Echo {
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
