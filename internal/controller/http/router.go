// Package v1 implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"go-eprescription-clean/config"
	_ "go-eprescription-clean/docs" // Swagger docs.
	"go-eprescription-clean/internal/controller/http/middleware"
	v1 "go-eprescription-clean/internal/controller/http/v1"
	"go-eprescription-clean/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go E-Presciprion Clean API
// @description Example of clean architecture in Go.
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, u v1.Usecases, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(middleware.CORS(l))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("my-service-name")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewSignaRoutes(apiV1Group, u.Signa, l)
		v1.NewPatientRoutes(apiV1Group, u.Patient, l)
		v1.NewMedicineRoutes(apiV1Group, u.Medicine, l)
		v1.NewTransactionRoutes(apiV1Group, u.Transaction, u.MedicineDetail, u.Medicine, l)
		v1.NewMedicineDetailRoutes(apiV1Group, u.MedicineDetail, l)
		// Add other routes here
	}
}
