package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/presenters"

	"github.com/labstack/echo/v4"
)

func (ctl *workerController) IngestForm(c echo.Context) error {
	ctx := c.Request().Context()

	req := &presenters.IngestedForm{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	form, err := req.ToModel()
	if err != nil {
		slog.Error("failed to fingerprint ingested form", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := form.Validate(); err != nil {
		slog.Error("invalid ingested form", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	created, ingestedForm, err := ctl.ingestedSrv.GetOrCreate(ctx, form)
	if err != nil {
		slog.Error("failed to create ingested form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Duplicate delivery of identical content: if already fully processed, acknowledge.
	if !created {
		transformed, tErr := ctl.transformedSrv.GetByIngestedFormID(ctx, ingestedForm.ID)
		if tErr == nil && transformed.SentToBot {
			return c.JSON(http.StatusOK, ingestedForm)
		}
		if tErr == nil && !transformed.SentToBot {
			if err := ctl.queueSrv.AddObject(ctx, "send-to-bot", &presenters.SendToBotRequest{
				TransformedID: transformed.ID,
			}); err != nil {
				slog.Error("failed to re-queue send-to-bot", "error", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusOK, ingestedForm)
		}
		if !errors.Is(tErr, models.ErrTransformedFormNotFound) && tErr != nil {
			slog.Error("failed to lookup transformed form", "error", tErr)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": tErr.Error()})
		}
	}

	if err := ctl.queueSrv.AddObject(ctx, "transform", &presenters.TransformRequest{
		IngestedID: ingestedForm.ID,
	}); err != nil {
		slog.Error("failed to add object to queue", "error", err)
		ingestedForm.Status = models.IngestedStatusFailed
		if _, patchErr := ctl.ingestedSrv.Patch(ctx, ingestedForm, "status"); patchErr != nil {
			slog.Error("failed to mark ingested form failed", "error", patchErr)
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ingestedForm)
}
