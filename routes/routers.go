package routes

import (
	"github.com/LitPad/backend/internetcomputer"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Endpoint struct {
	DB *gorm.DB
}

type WalletService struct {
	WS *internetcomputer.WalletService
	DB *gorm.DB
}

func SetupRoutes(app *fiber.App, db *gorm.DB, ws *internetcomputer.WalletService) {
	endpoint := Endpoint{DB: db}

	walletService := WalletService{
		DB: db,
		WS: ws,
	}

	// ROUTES (40)
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

	// Profile Routes (6)
	profilesRouter := api.Group("/profiles")
	profilesRouter.Get("/profile/:username", endpoint.GetProfile)
	profilesRouter.Patch("/update", endpoint.AuthMiddleware, endpoint.UpdateProfile)
	profilesRouter.Put("/update-password", endpoint.AuthMiddleware, endpoint.UpdatePassword)
	profilesRouter.Get("/profile/:username/follow", endpoint.AuthMiddleware, endpoint.FollowUser)
	profilesRouter.Get("/notifications", endpoint.AuthMiddleware, endpoint.GetNotifications)
	profilesRouter.Post("/notifications/read", endpoint.AuthMiddleware, endpoint.ReadNotification)

	// Book Routes (24)
	bookRouter := api.Group("/books")
	bookRouter.Get("", endpoint.GetLatestBooks)
	bookRouter.Post("", endpoint.AuthorMiddleware, endpoint.CreateBook)
	bookRouter.Get("/bought", endpoint.AuthMiddleware, endpoint.GetBoughtBooks)
	bookRouter.Get("/book/:slug", endpoint.GetSingleBook)
	bookRouter.Get("/book/:slug/chapters", endpoint.AuthOrGuestMiddleware, endpoint.GetBookChapters)
	bookRouter.Post("/book/:slug", endpoint.AuthMiddleware, endpoint.ReviewBook)
	bookRouter.Put("/book/review/:id", endpoint.AuthMiddleware, endpoint.EditBookReview)
	bookRouter.Delete("/book/review/:id", endpoint.AuthMiddleware, endpoint.DeleteBookReview)
	bookRouter.Get("/book/review/:id/replies", endpoint.GetReviewReplies)
	bookRouter.Post("/book/review/:id/replies", endpoint.AuthMiddleware, endpoint.ReplyReview)
	bookRouter.Put("/book/review/replies/:id", endpoint.AuthMiddleware, endpoint.EditReply)
	bookRouter.Delete("/book/review/replies/:id", endpoint.AuthMiddleware, endpoint.DeleteReply)
	bookRouter.Get("/book/:slug/vote", endpoint.AuthMiddleware, endpoint.VoteBook)
	bookRouter.Get("/lanterns-generation/:amount", endpoint.AuthMiddleware, endpoint.ConvertCoinsToLanterns)

	bookRouter.Put("/book/:slug", endpoint.AuthorMiddleware, endpoint.UpdateBook)
	bookRouter.Delete("/book/:slug", endpoint.AuthorMiddleware, endpoint.DeleteBook)
	bookRouter.Get("/book/:slug/buy", endpoint.AuthMiddleware, endpoint.BuyABook)
	bookRouter.Get("/book/:slug/buy-chapter", endpoint.AuthMiddleware, endpoint.BuyAChapter)
	bookRouter.Post("/book/:slug/set-contract", endpoint.AuthorMiddleware, endpoint.SetContract)
	bookRouter.Put("/book/chapter/:slug", endpoint.AuthorMiddleware, endpoint.UpdateChapter)
	bookRouter.Delete("/book/chapter/:slug", endpoint.AuthorMiddleware, endpoint.DeleteChapter)
	bookRouter.Post("/book/:slug/add-chapter", endpoint.AuthorMiddleware, endpoint.AddChapter)
	bookRouter.Get("/author/:username", endpoint.GetLatestAuthorBooks)
	bookRouter.Get("/genres", endpoint.GetAllBookGenres)
	bookRouter.Get("/tags", endpoint.GetAllBookTags)

	// Gifts Routes (4)
	giftsRouter := api.Group("/gifts")
	giftsRouter.Get("", endpoint.GetAllGifts)
	giftsRouter.Get("/:username/:gift_slug/send", endpoint.AuthMiddleware, endpoint.SendGift)
	giftsRouter.Get("/sent", endpoint.AuthorMiddleware, endpoint.GetAllSentGifts)
	giftsRouter.Get("/sent/:id/claim", endpoint.AuthorMiddleware, endpoint.ClaimGift)

	// Wallet Routes (7)
	walletRouter := api.Group("/wallet")
	walletRouter.Get("/coins", endpoint.AvailableCoins)
	walletRouter.Post("/coins", endpoint.AuthMiddleware, endpoint.BuyCoins)
	walletRouter.Get("/transactions", endpoint.AuthMiddleware, endpoint.AllUserTransactions)
	walletRouter.Post("/verify-payment", endpoint.VerifyPayment)
	walletRouter.Get("/plans", endpoint.GetSubscriptionPlans)
	walletRouter.Put("/plans", endpoint.AdminMiddleware, endpoint.UpdateSubscriptionPlan)
	walletRouter.Post("/subscription", endpoint.AuthMiddleware, endpoint.BookSubscription)

	// Internet Computer
	walletRouter.Get("balance", walletService.GetOnChainBalance)

	// ADMIN ROUTES (7)
	adminRouter := api.Group("/admin")
	// Admin Users
	adminRouter.Put("/", endpoint.AdminMiddleware, endpoint.UpdateProfile)
	adminRouter.Get("/users", endpoint.AdminMiddleware, endpoint.AdminGetUsers)
	adminRouter.Put("/users/user", endpoint.AdminMiddleware, endpoint.AdminUpdateUser)
	adminRouter.Get("/users/user/:username/toggle-activation", endpoint.AdminMiddleware, endpoint.ToggleUserActivation)

	// Admin Books (2)
	adminRouter.Get("/books", endpoint.AdminMiddleware, endpoint.AdminGetBooks)
	adminRouter.Get("/contracts", endpoint.AdminMiddleware, endpoint.AdminGetBookContracts)

	// Admin Waitlist (1)
	adminRouter.Get("/waitlist", endpoint.AdminMiddleware, endpoint.AdminGetWaitlist)

	// Admin Payments (1)
	adminRouter.Get("/payments/transactions", endpoint.AdminMiddleware, endpoint.AdminGetTransactions)
	// --------------------------------------------------------------------------------

	// Waitlist Routes (1)
	api.Post("/waitlist", endpoint.AddToWaitlist)

	// Register Sockets (1)
	api.Get("/ws/notifications", websocket.New(endpoint.NotificationSocket))
}
