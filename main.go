package main

import (
	"log"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/database"
	_ "github.com/LitPad/backend/docs"
	"github.com/LitPad/backend/initials"
	"github.com/LitPad/backend/jobs"
	"github.com/LitPad/backend/routes"
	"github.com/LitPad/backend/templates"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html/v2"
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
	// Load config
	conf := config.GetConfig()

	// Get Database
	db := database.ConnectDb(conf)
	// Create initial data
	initials.CreateInitialData(db, conf)

	engine := html.New("./templates", ".html")
	engine.AddFunc("add", templates.TemplateFuncMap["add"])
	engine.AddFunc("sub", templates.TemplateFuncMap["sub"])
	engine.AddFunc("sequence", templates.TemplateFuncMap["sequence"])
	engine.AddFunc("paginationRange", templates.TemplateFuncMap["paginationRange"])

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     conf.CORSAllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin, Content-Disposition",
		AllowCredentials: conf.CORSAllowCredentials,
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// General Rate Limiter: 200 requests per 1 minute
	app.Use(limiter.New(limiter.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	return c.IP() == "127.0.0.1" // Do not set a rate limiter for local requests
		// },
		Max:        100,                         // maximum number of requests
		Expiration: 1 * time.Minute,             // time period for limiting (1 minute)
		KeyGenerator: func(c *fiber.Ctx) string { // define how to generate the key
			return c.IP() // Limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(utils.RateLimitError("Rate Limit Reached"))
		},
	}))

	// Swagger Config
	swaggerCfg := swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "/",
		Title:    "LITPAD API Documentation",
		CacheAge: 1,
	}

	app.Use(swagger.New(swaggerCfg))
	app.Use(routes.RequestLogger(db))

	// Register Routes & Sockets
	routes.SetupRoutes(app, db)
	
	// RUN JOBS
	jobs.RunJobs(conf, db)
	log.Fatal(app.Listen(":" + conf.Port)) 

}