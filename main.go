package main

import (
	"log"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/database"
	_ "github.com/LitPad/backend/docs"
	"github.com/LitPad/backend/initials"
	"github.com/LitPad/backend/jobs"
	"github.com/LitPad/backend/routes"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title LITPAD API
// @description.markdown api
// @version 1.0
// @Accept json
// @Produce json
// @BasePath  /api/v1
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer jwt_string' to correctly set the API Key
func main() {
	log.Fatal("Simulating failure for testing purposes")

	// Load config
	conf := config.GetConfig()

	// Get Database
	db := database.ConnectDb(conf)
	// Create initial data
	initials.CreateInitialData(db, conf)

	app := fiber.New()

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     conf.CORSAllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin, Content-Disposition",
		AllowCredentials: conf.CORSAllowCredentials,
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
	routes.SetupRoutes(app, db)

	// RUN JOBS
	jobs.RunJobs(conf, db)
	log.Fatal(app.Listen(":" + conf.Port))

}
