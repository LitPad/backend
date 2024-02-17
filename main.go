package main

import (
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/LitPad/backend/docs"
)

// @title LITPAD API
// @version 4.0
// @Accept json
// @Produce json
// @BasePath  /api/v1
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer jwt_string' to correctly set the API Key
func main() {
	app := fiber.New()

	// CORS config
	app.Use(cors.New(cors.Config{
		// AllowOrigins:     ,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin, Content-Disposition",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Swagger Config
	swaggerCfg := swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "/",
		Title:    "LITPAD API Documentation",
		CacheAge: 1,
	}

	app.Use(swagger.New(swaggerCfg))
	// Register Routes & Sockets
	// routes.SetupRoutes(app, db)
	log.Fatal(app.Listen(":8000"))
}