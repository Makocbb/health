package controllers

import (
	"go-take-home-test/internal/presenters"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ctl *workerController) IngestForm(c echo.Context) error {
	ctx := c.Request().Context()

	req := &presenters.IngestedForm{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ingestedForm, err := ctl.ingestedSrv.Create(ctx, req.ToModel())
	if err != nil {
		slog.Error("failed to create ingested form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := ctl.queueSrv.AddObject(ctx, "transform", &presenters.TransformRequest{
		IngestedID: ingestedForm.ID,
	}); err != nil {
		slog.Error("failed to add object to queue", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ingestedForm)
}
