package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"go-take-home-test/internal/models"
	"go-take-home-test/internal/presenters"

	"github.com/labstack/echo/v4"
)

// RetryForm re-queues a failed/incomplete form for transform or send-to-bot.
func (ctl *workerController) RetryForm(c echo.Context) error {
	ctx := c.Request().Context()

	req := &presenters.RetryRequest{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var ingested *models.IngestedForm
	var err error
	switch {
	case req.IngestedID != 0:
		ingested, err = ctl.ingestedSrv.GetByID(ctx, req.IngestedID)
	case req.SessionID != "":
		ingested, err = ctl.ingestedSrv.GetBySessionID(ctx, req.SessionID)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ingested_id or session_id is required"})
	}
	if err != nil {
		slog.Error("failed to get ingested form for retry", "error", err)
		status := http.StatusInternalServerError
		if errors.Is(err, models.ErrIngestedFormNotFound) {
			status = http.StatusNotFound
		}
		return c.JSON(status, map[string]string{"error": err.Error()})
	}

	transformed, tErr := ctl.transformedSrv.GetByIngestedFormID(ctx, ingested.ID)
	if tErr == nil {
		if transformed.SentToBot {
			return c.JSON(http.StatusOK, map[string]any{
				"status":  "already_sent",
				"message": "form was already delivered to FORM-BOT",
			})
		}
		if err := ctl.queueSrv.AddObject(ctx, "send-to-bot", &presenters.SendToBotRequest{
			TransformedID: transformed.ID,
		}); err != nil {
			slog.Error("failed to re-queue send-to-bot", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]any{
			"status":          "requeued_send_to_bot",
			"transformed_id":  transformed.ID,
			"ingested_id":     ingested.ID,
		})
	}
	if !errors.Is(tErr, models.ErrTransformedFormNotFound) {
		slog.Error("failed to lookup transformed form for retry", "error", tErr)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": tErr.Error()})
	}

	ingested.Status = models.IngestedStatusPending
	if _, err := ctl.ingestedSrv.Patch(ctx, ingested, "status"); err != nil {
		slog.Error("failed to reset ingested status", "error", err)
	}

	if err := ctl.queueSrv.AddObject(ctx, "transform", &presenters.TransformRequest{
		IngestedID: ingested.ID,
	}); err != nil {
		slog.Error("failed to re-queue transform", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":      "requeued_transform",
		"ingested_id": ingested.ID,
	})
}
