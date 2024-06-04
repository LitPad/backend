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

	// Auth Routes (11)
	authRouter := api.Group("/auth")
	authRouter.Post("/register", endpoint.Register)
	authRouter.Post("/verify-email", endpoint.VerifyEmail)
	authRouter.Post("/resend-verification-email", endpoint.ResendVerificationEmail)
	authRouter.Post("/send-password-reset-otp", endpoint.SendPasswordResetOtp)
	authRouter.Get("/verify-password-reset-token/:token_string", endpoint.VerifyPasswordResetToken)
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

	// Book Routes ()
	bookRouter := api.Group("/books")
	bookRouter.Get("", endpoint.GetLatestBooks)
	bookRouter.Post("", endpoint.AuthorMiddleware, endpoint.CreateBook)
	bookRouter.Get("/bought", endpoint.AuthMiddleware, endpoint.GetBoughtBooks)
	bookRouter.Put("/book/:slug", endpoint.AuthorMiddleware, endpoint.UpdateBook)
	bookRouter.Delete("/book/:slug", endpoint.AuthorMiddleware, endpoint.DeleteBook)
	bookRouter.Get("/book/:slug/buy", endpoint.AuthMiddleware, endpoint.BuyBook)
	bookRouter.Put("/book/chapter/:slug", endpoint.AuthorMiddleware, endpoint.UpdateChapter)
	bookRouter.Delete("/book/chapter/:slug", endpoint.AuthorMiddleware, endpoint.DeleteChapter)
	bookRouter.Post("/book/:slug/add-chapter", endpoint.AuthorMiddleware, endpoint.AddChapter)
	bookRouter.Get("/author/:username", endpoint.GetLatestAuthorBooks)
	bookRouter.Get("/genres", endpoint.GetAllBookGenres)
	bookRouter.Get("/tags", endpoint.GetAllBookTags)

	// Wallet Routes ()
	walletRouter := api.Group("/wallet")
	walletRouter.Get("/coins", endpoint.AvailableCoins)
	walletRouter.Post("/coins", endpoint.AuthMiddleware, endpoint.BuyCoins)
	walletRouter.Get("/transactions", endpoint.AuthMiddleware, endpoint.AllUserTransactions)
	walletRouter.Post("/verify-payment", endpoint.VerifyPayment)

	// Admin Routes ()
	adminRouter := api.Group("/admin")
	adminRouter.Get("/users", endpoint.AuthMiddleware, endpoint.GetUsers)
	adminRouter.Get("/books", endpoint.AuthMiddleware, endpoint.GetBooks)
}
