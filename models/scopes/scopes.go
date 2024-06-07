package scopes

import (
	"github.com/LitPad/backend/models"
	"gorm.io/gorm"
)

func VerifiedUserScope(db *gorm.DB) *gorm.DB {
	return db.Where(models.User{IsEmailVerified: true})
}

func FollowerFollowingPreloaderScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(VerifiedUserScope).Preload("Followers").Preload("Followers.Followers").Preload("Followings").Preload("Followings.Followers").Preload("Followings.Books")
}

func AuthorGenreTagBookScope(db *gorm.DB) *gorm.DB {
	return db.Joins("Author").Joins("Genre").Preload("Tags").Preload("Chapters")
}

func AuthorGenreTagReviewsBookScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(AuthorGenreTagBookScope).Preload("Reviews").Preload("Reviews.User").Preload("Reviews.Likes").Preload("Reviews.Replies")
}

func BoughtAuthorGenreTagBookScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Book").Preload("Book.Author").Preload("Book.Genre").Preload("Book.Tags").Preload("Book.Chapters")
}

func BoughtAuthorGenreTagReviewsBookScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(BoughtAuthorGenreTagBookScope).Preload("Book.Reviews").Preload("Book.Reviews.User").Preload("Book.Reviews.Likes").Preload("Book.Reviews.Replies")
}

func SentGiftRelatedScope(db *gorm.DB) *gorm.DB {
	return db.Joins("Sender").Joins("Receiver").Joins("Gift")
}