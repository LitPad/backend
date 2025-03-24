package tests

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/routes"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const MASTER_PASSWORD = "testpassword"

// AUTH FIXTURES
func TestUser(db *gorm.DB) models.User {
	email := "testuser@example.com"
	db.Where("email = ?", email).Delete(&models.User{})

	user := models.User{
		Email:    email,
		Password: MASTER_PASSWORD,
	}
	db.Create(&user)
	return user
}

func TestVerifiedUser(db *gorm.DB, activeSub ...bool) models.User {
	email := "testverifieduser@example.com"

	user := models.User{
		Email:           email,
		Password:        MASTER_PASSWORD,
		IsEmailVerified: true,
	}
	if len(activeSub) > 0 {
		user.Email = "testactivesubuser@example.com"
		expiry := time.Now().AddDate(0, 1, 0)
		user.SubscriptionExpiry = &expiry
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func TestAuthor(db *gorm.DB, another ...bool) models.User {
	email := "testauthormail@example.com"
	user := models.User{
		Email:           email,
		Password:        MASTER_PASSWORD,
		IsEmailVerified: true,
		AccountType:     choices.ACCTYPE_AUTHOR,
	}

	if len(another) > 0 {
		user.Email = "testanotherauthormail@example.com"
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func JwtData(db *gorm.DB, user models.User) models.User {
	access := routes.GenerateAccessToken(user)
	refresh := routes.GenerateRefreshToken()
	user.Access = &access
	user.Refresh = &refresh
	db.Save(&user)
	return user
}

func AccessToken(db *gorm.DB, user models.User) string {
	user = JwtData(db, user)
	return *user.Access
}

// BOOKS TEST DATA
func TagData(db *gorm.DB) models.Tag {
	tag := models.Tag{Name: "Test Tag"}
	db.FirstOrCreate(&tag, tag)
	return tag
}

func GenreData(db *gorm.DB) models.Genre {
	tag := TagData(db)
	genre := models.Genre{Name: "Test Genre", Tags: []models.Tag{tag}}
	db.Omit("Tags.*").FirstOrCreate(&genre, models.Genre{Name: "Test Genre"})
	genre.Tags = []models.Tag{tag}
	return genre
}

func BookData(db *gorm.DB, user models.User) models.Book {
	book := models.Book{
		AuthorID: user.ID, Title: "Test Book", Blurb: "blurning me",
		AgeDiscretion: choices.ATYPE_EIGHTEEN, GenreID: GenreData(db).ID,
		CoverImage: "https://coverimage.url", Tags: []models.Tag{TagData(db)},
	}
	db.Omit("Tags.*").FirstOrCreate(&book, book)
	return book
}

func ChapterData(db *gorm.DB, book models.Book) models.Chapter {
	chapter := models.Chapter{BookID: book.ID, Title: "Test Chapter", Text: "Stop doing that"}
	db.FirstOrCreate(&chapter, chapter)
	return chapter
}

func ReviewData(db *gorm.DB, book models.Book, user models.User) models.Review {
	review := models.Review{BookID: book.ID, UserID: user.ID, Rating: choices.RC_1, Text: "This is a test review"}
	db.FirstOrCreate(&review, models.Review{BookID: book.ID, UserID: user.ID})
	return review
}

func ReplyData(db *gorm.DB, review models.Review, user models.User) models.Reply {
	reply := models.Reply{ReviewID: &review.ID, UserID: user.ID, Text: "This is a test reply"}
	db.FirstOrCreate(&reply, reply)
	return reply
}

func ParagraphCommentData(db *gorm.DB, chapter models.Chapter, user models.User) models.ParagraphComment {
	comment := models.ParagraphComment{ChapterID: chapter.ID, Index: 2, UserID: user.ID, Text: "This is a test comment"}
	db.FirstOrCreate(&comment, comment)
	return comment
}

// PROFILES TEST DATA

func NotificationData(db *gorm.DB, sender models.User, receiver models.User) models.Notification {
	notification := models.Notification{
		SenderID: sender.ID, ReceiverID: receiver.ID,
		Ntype: choices.NT_FOLLOWING, Text: "This is a test notification",
	}
	db.FirstOrCreate(&notification, notification)
	return notification
}

// ADMIN TEST DATA
func TestAdmin(db *gorm.DB) models.User {
	email := "testadmin@example.com"

	user := models.User{
		Email:           email,
		Password:        MASTER_PASSWORD,
		IsStaff:         true,
		IsSuperuser:     true,
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}


func TestSubscriber(db *gorm.DB) models.User {
	email := "testsubscriber@example.com"
	subExpiry := time.Now().AddDate(0, 1, 0)
	currentPlan := choices.ST_MONTHLY
	user := models.User{
		Email:           email,
		Password:        MASTER_PASSWORD,
		IsEmailVerified: true,
		SubscriptionExpiry: &subExpiry,
		CurrentPlan: &currentPlan,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func TestSubscriptionPlan(db *gorm.DB) models.SubscriptionPlan {
	plan := models.SubscriptionPlan{
		Amount: decimal.NewFromInt(1000), SubType: choices.ST_MONTHLY,
	}
	db.FirstOrCreate(&plan, plan)
	return plan
}

func TestTransaction(db *gorm.DB, user models.User) models.Transaction {
	plan := TestSubscriptionPlan(db)
	transaction := models.Transaction{
		UserID: user.ID, SubscriptionPlanID: &plan.ID,
		PaymentType: choices.PTYPE_GPAY, PaymentPurpose: choices.PP_SUB,
	}
	db.FirstOrCreate(&transaction, transaction)
	return transaction
}

// WALLET TEST DATA
func TestCoin(db *gorm.DB) models.Coin {
	coin := models.Coin{Amount: 100, Price: decimal.NewFromInt(100)}
	db.FirstOrCreate(&coin, coin)
	return coin
}