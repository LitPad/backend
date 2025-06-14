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
	profilesRouter := api.Group("/profiles", endpoint.AuthMiddleware)
	profilesRouter.Get("/profile/:username", endpoint.GetProfile)
	profilesRouter.Patch("/update", endpoint.UpdateProfile)
	profilesRouter.Put("/update-password", endpoint.UpdatePassword)
	profilesRouter.Get("/profile/:username/follow", endpoint.FollowUser)
	profilesRouter.Get("/notifications", endpoint.GetNotifications)
	profilesRouter.Post("/notifications/read", endpoint.ReadNotification)

	// Book Routes (24)
	bookRouter := api.Group("/books")
	bookRouter.Get("", endpoint.GetLatestBooks)
	bookRouter.Post("", endpoint.AdminMiddleware, endpoint.CreateBook)
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

	bookRouter.Put("/book/:slug", endpoint.AdminMiddleware, endpoint.UpdateBook)
	bookRouter.Delete("/book/:slug", endpoint.AdminMiddleware, endpoint.DeleteBook)
	bookRouter.Post("/book/:slug/set-contract", endpoint.AdminMiddleware, endpoint.SetContract)
	bookRouter.Put("/book/chapter/:slug", endpoint.AdminMiddleware, endpoint.UpdateChapter)
	bookRouter.Delete("/book/chapter/:slug", endpoint.AdminMiddleware, endpoint.DeleteChapter)
	bookRouter.Post("/book/:slug/add-chapter", endpoint.AdminMiddleware, endpoint.AddChapter)

	bookRouter.Get("/book/chapters/chapter/:slug", endpoint.AuthMiddleware, endpoint.GetBookChapter)
	bookRouter.Get("/book/chapters/chapter/:slug/paragraph/:index/comments", endpoint.AuthMiddleware, endpoint.GetParagraphComments)
	bookRouter.Post("/book/chapters/chapter/:slug/paragraph/:index/comments", endpoint.AuthMiddleware, endpoint.AddParagraphComment)
	bookRouter.Put("/book/chapters/chapter/paragraph-comment/:id", endpoint.AuthMiddleware, endpoint.EditParagraphComment)
	bookRouter.Delete("/book/chapters/chapter/paragraph-comment/:id", endpoint.AuthMiddleware, endpoint.DeleteParagraphComment)

	bookRouter.Get("/book/chapters/chapter/comment/:id", endpoint.AuthMiddleware, endpoint.LikeAComment)

	bookRouter.Get("/author/:username", endpoint.GetLatestAuthorBooks)
	bookRouter.Get("/genres", endpoint.GetAllBookGenres)
	bookRouter.Get("/sections", endpoint.GetAllBookSections)
	bookRouter.Get("/sub-sections", endpoint.GetAllBookSubSections)
	bookRouter.Get("/tags", endpoint.GetAllBookTags)

	// Gifts Routes (4)
	giftsRouter := api.Group("/gifts")
	giftsRouter.Get("", endpoint.GetAllGifts)
	giftsRouter.Get("/:username/:gift_slug/send", endpoint.AuthMiddleware, endpoint.SendGift)
	giftsRouter.Get("/sent", endpoint.GetAllSentGifts)
	giftsRouter.Get("/sent/:id/claim", endpoint.ClaimGift)

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
	adminRouter := api.Group("/admin", endpoint.AdminMiddleware)
	adminRouter.Get("/", endpoint.AdminDashboard)
	adminUsersRouter := adminRouter.Group("/users", endpoint.AdminMiddleware)
	// Admin Users
	adminRouter.Put("/", endpoint.UpdateProfile)
	adminUsersRouter.Get("", endpoint.AdminGetUsers)
	adminUsersRouter.Put("/:username", endpoint.AdminUpdateUser)
	adminUsersRouter.Get("/:username/toggle-activation", endpoint.ToggleUserActivation)
	adminUsersRouter.Post("/admins/invite", endpoint.InviteAdmin)

	adminRouter.Get("/subscribers", endpoint.AdminGetSubscribers)

	// Admin Books (2)
	adminBooksRouter := adminRouter.Group("/books", endpoint.AdminMiddleware)
	adminBooksRouter.Get("", endpoint.AdminGetBooks)
	adminBooksRouter.Get("/by-username/:username", endpoint.AdminGetAuthorBooks)
	adminBooksRouter.Get("/book-detail/:slug", endpoint.AdminGetBookDetails)
	adminBooksRouter.Get("/book-detail/:slug/reading-progress", endpoint.AdminGetBookReadingProgress)
	adminBooksRouter.Get("/book-detail/:slug/retention-stats", endpoint.AdminGetBookRetentionStats)
	adminBooksRouter.Get("/contracts", endpoint.AdminGetBookContracts)
	adminBooksRouter.Post("/genres", endpoint.AdminAddBookGenre)
	adminBooksRouter.Post("/tags/add/:genre_slug", endpoint.AdminAddBookTag)
	adminBooksRouter.Get("/sections", endpoint.AdminGetSections)
	adminBooksRouter.Post("/sections", endpoint.AdminAddBookSection)
	adminBooksRouter.Post("/sections/:slug/subsections", endpoint.AdminAddBookSubSection)
	adminBooksRouter.Put("/genres/:slug", endpoint.AdminUpdateBookGenre)
	adminBooksRouter.Put("/sections/:slug", endpoint.AdminUpdateBookSection)
	adminBooksRouter.Get("/subsections/:slug", endpoint.AdminGetSubSection)
	adminBooksRouter.Put("/subsections/:slug", endpoint.AdminUpdateBookSubSection)
	adminBooksRouter.Put("/tags/:slug", endpoint.AdminUpdateBookTag)
	adminBooksRouter.Delete("/genres/:slug", endpoint.AdminDeleteBookGenre)
	adminBooksRouter.Delete("/sections/:slug", endpoint.AdminDeleteBookSection)
	adminBooksRouter.Delete("/subsections/:slug", endpoint.AdminDeleteBookSubSection)
	adminBooksRouter.Get("/subsections/:slug/add-book/:book_slug", endpoint.AddBookToSubSection)
	adminBooksRouter.Get("/subsections/:slug/remove-book/:book_slug", endpoint.RemoveBookFromSubSection)
	adminBooksRouter.Get("/book/:slug/toggle-book-completion-status", endpoint.ToggleBookCompletionStatus)
	adminBooksRouter.Delete("/tags/:slug", endpoint.AdminDeleteBookTag)

	// Admin Contents
	adminRouter.Get("/featured-contents", endpoint.AdminGetFeaturedContents)
	adminRouter.Post("/featured-contents", endpoint.AdminAddAFeaturedContent)
	adminRouter.Put("/featured-contents/:id", endpoint.AdminUpdateAFeaturedContent)
	adminRouter.Delete("/featured-contents/:id", endpoint.AdminDeleteAFeaturedContent)

	// Admin Waitlist (1)
	adminRouter.Get("/waitlist", endpoint.AdminGetWaitlist)

	// Admin Payments (1)
	walletRouter.Put("/payments/plans", endpoint.UpdateSubscriptionPlan)
	adminRouter.Get("/payments/transactions", endpoint.AdminGetTransactions)
	// --------------------------------------------------------------------------------

	// Waitlist Routes (1)
	api.Post("/waitlist", endpoint.AddToWaitlist)

	// Register Sockets (1)
	api.Get("/ws/notifications", websocket.New(endpoint.NotificationSocket))
}
