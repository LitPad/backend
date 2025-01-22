package tests

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/routes"
	"gorm.io/gorm"
)

// AUTH FIXTURES
func TestUser(db *gorm.DB) models.User {
	email := "testuser@example.com"
	db.Where("email = ?", email).Delete(&models.User{})

	user := models.User{
		Email:          email,
		Password:       "testpassword",
	}
	db.Create(&user)
	return user
}

func TestVerifiedUser(db *gorm.DB) models.User {
	email := "testverifieduser@example.com"
	db.Where("email = ?", email).Delete(&models.User{})

	user := models.User{
		Email:          email,
		Password:       "testpassword",
		IsEmailVerified: true,
	}
	db.Create(&user)
	return user
}

func TestAuthor(db *gorm.DB, another ...bool) models.User {
	email := "testauthormail@example.com"
	user := models.User{
		Email:          email,
		Password:       "testpassword",
		IsEmailVerified: true,
		AccountType: choices.ACCTYPE_AUTHOR,
	}

	if len(another) > 0 {
		user.Email = "testanotherauthormail@example.com"
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}


func JwtData(db *gorm.DB, user models.User) models.User {
	access := routes.GenerateAccessToken(user.ID)
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