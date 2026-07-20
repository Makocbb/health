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

// Routes registers route handlers for the health service
func (ctl *workerController) Routes(g *echo.Group) {
	g.POST("/ingest", ctl.IngestForm)
	g.POST("/transform", ctl.TransformIngestedForm)
	g.POST("/send-to-bot", ctl.SendToBot)
}
