package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Endpoint struct {
	DB *gorm.DB
}

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	endpoint := Endpoint{DB: db}

	api := app.Group("/api/v1")

	// HealthCheck Route
	api.Get("/healthcheck", HealthCheck)

	// General Routes
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", endpoint.GetSiteDetails)
	generalRouter.Post("/subscribe", endpoint.Subscribe)
}
