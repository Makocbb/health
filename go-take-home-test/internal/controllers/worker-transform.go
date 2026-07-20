package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/presenters"

	"github.com/labstack/echo/v4"
)

func (ctl *workerController) TransformIngestedForm(c echo.Context) error {
	ctx := c.Request().Context()

	req := &presenters.TransformRequest{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ingestedForm, err := ctl.ingestedSrv.GetByID(ctx, req.IngestedID)
	if err != nil {
		slog.Error("failed to get ingested form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Idempotent: reuse existing transformed row instead of creating duplicates.
	if existing, err := ctl.transformedSrv.GetByIngestedFormID(ctx, ingestedForm.ID); err == nil {
		if !existing.SentToBot {
			if err := ctl.queueSrv.AddObject(ctx, "send-to-bot", &presenters.SendToBotRequest{
				TransformedID: existing.ID,
			}); err != nil {
				slog.Error("failed to add object to queue", "error", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}
		return c.JSON(http.StatusOK, existing)
	} else if !errors.Is(err, models.ErrTransformedFormNotFound) {
		slog.Error("failed to lookup transformed form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	transformedForm, err := ingestedForm.ToTransformedForm()
	if err != nil {
		slog.Error("failed to transform ingested form", "error", err)
		ingestedForm.Status = models.IngestedStatusFailed
		if _, patchErr := ctl.ingestedSrv.Patch(ctx, ingestedForm, "status"); patchErr != nil {
			slog.Error("failed to mark ingested form failed", "error", patchErr)
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	coords, err := ctl.postcodeSrv.LookupPostcode(ctx, transformedForm.Postcode)
	if err != nil {
		slog.Error("failed to get coordinates", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	transformedForm.Longitude = coords.Longitude
	transformedForm.Latitude = coords.Latitude

	transformedForm, err = ctl.transformedSrv.Create(ctx, transformedForm)
	if err != nil {
		// Unique race: another worker created it first.
		if existing, findErr := ctl.transformedSrv.GetByIngestedFormID(ctx, ingestedForm.ID); findErr == nil {
			return c.JSON(http.StatusOK, existing)
		}
		slog.Error("failed to create transformed form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	ingestedForm.Status = models.IngestedStatusTransformed
	if _, err := ctl.ingestedSrv.Patch(ctx, ingestedForm, "status"); err != nil {
		slog.Error("failed to mark ingested form transformed", "error", err)
	}

	if err := ctl.queueSrv.AddObject(ctx, "send-to-bot", &presenters.SendToBotRequest{
		TransformedID: transformedForm.ID,
	}); err != nil {
		slog.Error("failed to add object to queue", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transformedForm)
}
