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

	// HealthCheck Route (1)
	api.Get("/healthcheck", HealthCheck)

	// General Routes (2)
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", endpoint.GetSiteDetails)
	generalRouter.Post("/subscribe", endpoint.Subscribe)

	// Auth Routes (10)
	authRouter := api.Group("/auth")
	authRouter.Post("/register", endpoint.Register)
	authRouter.Post("/verify-email", endpoint.VerifyEmail)
	authRouter.Post("/resend-verification-email", endpoint.ResendVerificationEmail)
	authRouter.Post("/send-password-reset-otp", endpoint.SendPasswordResetOtp)
	authRouter.Post("/set-new-password", endpoint.SetNewPassword)
	authRouter.Post("/login", endpoint.Login)
	authRouter.Post("/google", endpoint.GoogleLogin)
	authRouter.Post("/facebook", endpoint.FacebookLogin)
	authRouter.Post("/refresh", endpoint.Refresh)
	authRouter.Get("/logout", endpoint.AuthMiddleware, endpoint.Logout)

	// Profile Routes ()
	profilesRouter := api.Group("/profiles")
	profilesRouter.Get("/profile/:username", endpoint.GetProfile)
	profilesRouter.Patch("/update", endpoint.AuthMiddleware, endpoint.UpdateProfile)
	profilesRouter.Put("/update-password", endpoint.AuthMiddleware, endpoint.UpdatePassword)
	profilesRouter.Get("/profile/:username/follow", endpoint.AuthMiddleware, endpoint.FollowUser)

	// Admin Routes ()
	adminRouter := api.Group("/admin")
	adminRouter.Get("/users", endpoint.AuthMiddleware, endpoint.GetUsers)

}
