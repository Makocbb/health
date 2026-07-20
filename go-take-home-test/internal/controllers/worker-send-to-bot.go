package controllers

import (
	"encoding/json"
	"go-take-home-test/internal/models"
	"go-take-home-test/internal/presenters"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ctl *workerController) SendToBot(c echo.Context) error {
	ctx := c.Request().Context()

	req := &presenters.SendToBotRequest{}
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	transformedForm, err := ctl.transformedSrv.GetByID(ctx, req.TransformedID)
	if err != nil {
		slog.Error("failed to get transformed form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	body, err := json.Marshal(presenters.TransformedFromModel(transformedForm))
	if err != nil {
		slog.Error("failed to marshal transformed form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	transformLog, err := ctl.transformLogSrv.Create(ctx, &models.TransformLog{
		TransformedFormID: transformedForm.ID,
		Status:            models.TransformStatusPending,
	})
	if err != nil {
		slog.Error("failed to create transform log", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = ctl.emailSrv.SendEmail(ctx, models.EmailRequest{
		From:    "source@mock.com",
		To:      "destination@mock.com",
		Subject: "form-bot",
		Body:    body,
	})
	if err == nil {
		transformLog.Status = models.TransformStatusSuccess
		transformLog.Success = true
		if _, err := ctl.transformLogSrv.Patch(ctx, transformLog, "status", "success"); err != nil {
			slog.Error("failed to patch transform log", "error", err)
		}
		transformedForm.SentToBot = true
		if _, err := ctl.transformedSrv.Patch(ctx, transformedForm, "sent_to_bot"); err != nil {
			slog.Error("failed to patch transformed form", "error", err)
		}
	} else {
		transformLog.Status = models.TransformStatusFailed
		transformLog.Error = err.Error()
		if _, logError := ctl.transformLogSrv.Patch(ctx, transformLog, "status", "error", "message"); logError != nil {
			slog.Error("failed to patch transform log", "error", logError)
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{})
}
