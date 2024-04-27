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