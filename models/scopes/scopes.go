package scopes

import (
	"github.com/LitPad/backend/models"
	"gorm.io/gorm"
)

func VerifiedUserScope(db *gorm.DB) *gorm.DB {
	return db.Where(models.User{IsEmailVerified: true})
}

func FollowerFollowingPreloaderScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(VerifiedUserScope).Preload("Books").Preload("Followers").Preload("Followers.Followers").Preload("Followings").Preload("Followings.Followers").Preload("Followings.Books")
}

func FollowerFollowingUnVerifiedPreloaderScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Followers").Preload("Followers.Followers").Preload("Followings").Preload("Followings.Followers").Preload("Followings.Books")
}

func FollowerFollowingBooksPreloaderScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(FollowerFollowingPreloaderScope).Preload("Books").Preload("Followers.Followers").Preload("Followings").Preload("Followings.Followers").Preload("Followings.Books")
}

func AuthorGenreTagBookScope(db *gorm.DB) *gorm.DB {
	return db.Joins("Author").Joins("Genre").Joins("SubSection").Joins("SubSection.Section").Preload("Tags").Preload("Chapters").Preload("Votes").Preload("Reads")
}

func AuthorGenreTagBookPreloadScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Author").Preload("Genre").Preload("SubSection").Preload("SubSection.Section").Preload("Tags").Preload("Chapters").Preload("Votes").Preload("Reads")
}

func TagsChaptersVotesBookScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Tags").Preload("Chapters").Preload("Votes").Preload("Reads")
}

func AuthorGenreTagReviewsBookScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(AuthorGenreTagBookPreloadScope).Preload("Reviews").Preload("Reviews.User").Preload("Reviews.Likes").Preload("Reviews.Replies")
}

func BoughtChapterScope(db *gorm.DB) *gorm.DB {
	return db.Joins("Chapter").Joins("Chapter.Book")
}

func SentGiftRelatedScope(db *gorm.DB) *gorm.DB {
	return db.Joins("Sender").Joins("Receiver").Joins("Gift")
}

func NotificationRelatedScope(db *gorm.DB) *gorm.DB {
	return db.Joins("Sender").Joins("Book")
}
