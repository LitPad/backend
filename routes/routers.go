package routes

import (
	"github.com/LitPad/backend/config"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type Endpoint struct {
	DB     *gorm.DB
	Config config.Config
	Store  *session.Store
}

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	store := session.New()
	endpoint := Endpoint{DB: db, Config: config.GetConfig(), Store: store}

	// ROUTES (40)
	api := app.Group("/api/v1")

	// HealthCheck Route (1)
	api.Get("/healthcheck", HealthCheck)

	// Logs Route (3)
	logsRouter := app.Group("/logs")
	logsRouter.Get("", endpoint.RenderLogs)
	logsRouter.Get("/login", endpoint.RenderLogsLogin)
	logsRouter.Post("/login", endpoint.HandleLogsLogin)
	logsRouter.Get("/logout", endpoint.HandleLogsLogout)

	// General Routes (2)
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", endpoint.GetSiteDetails)
	generalRouter.Post("/subscribe", endpoint.Subscribe)

	// Auth Routes (11)
	authRouter := api.Group("/auth")
	authRouter.Post("/register", endpoint.Register)
	authRouter.Post("/verify-email", endpoint.VerifyEmail)
	authRouter.Post("/resend-verification-email", endpoint.ResendVerificationEmail)
	authRouter.Post("/send-password-reset-link", endpoint.SendPasswordResetLink)
	authRouter.Get("/verify-password-reset-token/:token_string", endpoint.VerifyPasswordResetToken)
	authRouter.Post("/set-new-password", endpoint.SetNewPassword)
	authRouter.Post("/login", endpoint.Login)
	authRouter.Post("/google", endpoint.GoogleLogin)
	// authRouter.Post("/facebook", endpoint.FacebookLogin)
	authRouter.Post("/refresh", endpoint.Refresh)
	authRouter.Get("/logout", endpoint.AuthMiddleware, endpoint.Logout)
	authRouter.Get("/logout/all", endpoint.AuthMiddleware, endpoint.LogoutAll)

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
	bookRouter.Get("/bookmarked", endpoint.AuthMiddleware, endpoint.GetBookmarkedBooks)
	bookRouter.Get("/book/:slug/bookmark", endpoint.AuthMiddleware, endpoint.BookmarkBook)
	bookRouter.Post("/book/:slug/report", endpoint.AuthMiddleware, endpoint.ReportBook)
	bookRouter.Get("/book/:slug", endpoint.GetSingleBook)
	bookRouter.Get("/book/:slug/chapters", endpoint.AuthOrGuestMiddleware, endpoint.GetBookChapters)
	bookRouter.Post("/book/:slug", endpoint.AuthMiddleware, endpoint.ReviewBook)
	bookRouter.Put("/book/review/:id", endpoint.AuthMiddleware, endpoint.EditBookReview)
	bookRouter.Delete("/book/review/:id", endpoint.AuthMiddleware, endpoint.DeleteBookReview)
	bookRouter.Get("/book/review/:id/replies", endpoint.GetReviewReplies)
	bookRouter.Post("/book/review-or-paragraph-comment/:id/replies", endpoint.AuthMiddleware, endpoint.ReplyReviewOrParagraphComment)
	bookRouter.Put("/book/review-or-paragraph-comment/replies/:id", endpoint.AuthMiddleware, endpoint.EditReply)
	bookRouter.Delete("/book/review-or-paragraph-comment/replies/:id", endpoint.AuthMiddleware, endpoint.DeleteReply)
	bookRouter.Get("/book/:slug/vote", endpoint.AuthMiddleware, endpoint.VoteBook)
	bookRouter.Get("/lanterns-generation/:amount", endpoint.AuthMiddleware, endpoint.ConvertCoinsToLanterns)

	bookRouter.Put("/book/:slug", endpoint.AuthorMiddleware, endpoint.UpdateBook)
	bookRouter.Delete("/book/:slug", endpoint.AuthorMiddleware, endpoint.DeleteBook)
	bookRouter.Post("/book/:slug/set-contract", endpoint.AuthorMiddleware, endpoint.SetContract)
	bookRouter.Put("/book/chapter/:slug", endpoint.AuthorMiddleware, endpoint.UpdateChapter)
	bookRouter.Delete("/book/chapter/:slug", endpoint.AuthorMiddleware, endpoint.DeleteChapter)
	bookRouter.Post("/book/:slug/add-chapter", endpoint.AuthorMiddleware, endpoint.AddChapter)

	bookRouter.Get("/book/chapters/chapter/:slug", endpoint.AuthMiddleware, endpoint.GetBookChapter)
	bookRouter.Get("/book/chapters/chapter/:slug/paragraph/:index/comments", endpoint.AuthMiddleware, endpoint.GetParagraphComments)
	bookRouter.Post("/book/chapters/chapter/:slug/paragraph/:index/comments", endpoint.AuthMiddleware, endpoint.AddParagraphComment)
	bookRouter.Put("/book/chapters/chapter/paragraph-comment/:id", endpoint.AuthMiddleware, endpoint.EditParagraphComment)
	bookRouter.Delete("/book/chapters/chapter/paragraph-comment/:id", endpoint.AuthMiddleware, endpoint.DeleteParagraphComment)

	bookRouter.Get("/book/chapters/chapter/comment/:id", endpoint.AuthMiddleware, endpoint.LikeAComment)

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
	walletRouter.Post("/subscription", endpoint.AuthMiddleware, endpoint.BookSubscription)

	// ICP Wallet Routes
	icpWalletRouter := walletRouter.Group("/icp")
	icpWalletRouter.Post("/", endpoint.CreateICPWallet)
	icpWalletRouter.Get("/:username/balance", endpoint.GetICPWalletBalance)
	icpWalletRouter.Get("/gifts/:username/:gift_slug/send", endpoint.AuthMiddleware, endpoint.SendGiftViaICPWallet)

	// ADMIN ROUTES (7)
	adminRouter := api.Group("/admin")
	adminRouter.Get("/", endpoint.AdminMiddleware, endpoint.AdminDashboard)

	// Admin Users
	adminRouter.Put("/", endpoint.AdminMiddleware, endpoint.UpdateProfile)
	adminRouter.Get("/users", endpoint.AdminMiddleware, endpoint.AdminGetUsers)
	adminRouter.Put("/users/:username", endpoint.AdminMiddleware, endpoint.AdminUpdateUser)
	adminRouter.Get("/users/:username/toggle-activation", endpoint.AdminMiddleware, endpoint.ToggleUserActivation)
	adminRouter.Post("/users/admins/invite", endpoint.AdminMiddleware, endpoint.InviteAdmin)

	// Admin Users
	adminRouter.Get("/subscribers", endpoint.AdminMiddleware, endpoint.AdminGetSubscribers)

	// Admin Books (2)
	adminRouter.Get("/books", endpoint.AdminMiddleware, endpoint.AdminGetBooks)
	adminRouter.Get("/books/by-username/:username", endpoint.AdminMiddleware, endpoint.AdminGetAuthorBooks)
	adminRouter.Get("/books/book-detail/:slug", endpoint.AdminMiddleware, endpoint.AdminGetBookDetails)
	adminRouter.Get("/books/contracts", endpoint.AdminMiddleware, endpoint.AdminGetBookContracts)
	adminRouter.Post("/books/genres", endpoint.AdminMiddleware, endpoint.AdminAddBookGenre)
	adminRouter.Post("/books/tags", endpoint.AdminMiddleware, endpoint.AdminAddBookTag)
	adminRouter.Put("/books/genres/:slug", endpoint.AdminMiddleware, endpoint.AdminUpdateBookGenre)
	adminRouter.Put("/books/tags/:slug", endpoint.AdminMiddleware, endpoint.AdminUpdateBookTag)
	adminRouter.Delete("/books/genres/:slug", endpoint.AdminMiddleware, endpoint.AdminDeleteBookGenre)
	adminRouter.Delete("/books/tags/:slug", endpoint.AdminMiddleware, endpoint.AdminDeleteBookTag)

	// Admin Waitlist (1)
	adminRouter.Get("/waitlist", endpoint.AdminMiddleware, endpoint.AdminGetWaitlist)

	// Admin Payments (1)
	walletRouter.Put("/payments/plans", endpoint.AdminMiddleware, endpoint.UpdateSubscriptionPlan)
	adminRouter.Get("/payments/transactions", endpoint.AdminMiddleware, endpoint.AdminGetTransactions)
	// --------------------------------------------------------------------------------

	// Waitlist Routes (1)
	api.Post("/waitlist", endpoint.AddToWaitlist)

	// Register Sockets (1)
	api.Get("/ws/notifications", websocket.New(endpoint.NotificationSocket))
}
