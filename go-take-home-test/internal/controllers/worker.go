package controllers

import (
	"go-take-home-test/internal/services"

	"github.com/labstack/echo/v4"
)

type workerController struct {
	ingestedSrv     services.IngestedFormService
	transformedSrv  services.TransformedFormService
	transformLogSrv services.TransformLogService
	queueSrv        services.QueueService
	postcodeSrv     services.PostcodeService
	emailSrv        services.EmailService
}

// WorkerControllerInterface ...
type WorkerControllerInterface interface {
	Routes(g *echo.Group)
	PublicRoutes(e *echo.Echo)
}

// NewWorkerController ...
func NewWorkerController(
	ingestedSrv services.IngestedFormService,
	transformedSrv services.TransformedFormService,
	transformLogSrv services.TransformLogService,
	queueSrv services.QueueService,
	postcodeSrv services.PostcodeService,
	emailSrv services.EmailService,
) WorkerControllerInterface {
	return &workerController{
		ingestedSrv:     ingestedSrv,
		transformedSrv:  transformedSrv,
		transformLogSrv: transformLogSrv,
		queueSrv:        queueSrv,
		postcodeSrv:     postcodeSrv,
		emailSrv:        emailSrv,
	}
}

// Routes registers worker and public ingest/retry routes on the given group.
// Prefer mounting at "/workers" for internal steps; also register public routes on root.
func (ctl *workerController) Routes(g *echo.Group) {
	g.POST("/ingest", ctl.IngestForm)
	g.POST("/transform", ctl.TransformIngestedForm)
	g.POST("/send-to-bot", ctl.SendToBot)
	g.POST("/retry", ctl.RetryForm)
}

// PublicRoutes registers the README-facing /ingest and /retry endpoints.
func (ctl *workerController) PublicRoutes(e *echo.Echo) {
	e.POST("/ingest", ctl.IngestForm)
	e.POST("/retry", ctl.RetryForm)
}
