package controllers

import (
	"go-take-home-test/internal/presenters"
	"log/slog"
	"net/http"

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

	transformedForm, err := ingestedForm.ToTransformedForm()
	if err != nil {
		slog.Error("failed to transform ingested form", "error", err)
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
		slog.Error("failed to create transformed form", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := ctl.queueSrv.AddObject(ctx, "send-to-bot", &presenters.SendToBotRequest{
		TransformedID: transformedForm.ID,
	}); err != nil {
		slog.Error("failed to add object to queue", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transformedForm)
}
