package tests

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/routes"
	"gorm.io/gorm"
)

// AUTH FIXTURES
func TestUser(db *gorm.DB) models.User {
	user := models.User{
		Email:          "testuser@example.com",
		Password:       "testpassword",
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func TestVerifiedUser(db *gorm.DB) models.User {
	user := models.User{
		Email:           "testverifieduser@example.com",
		Password:        "testpassword",
		IsEmailVerified: true,
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

func AccessToken(db *gorm.DB) string {
	user := TestVerifiedUser(db)
	user = JwtData(db, user)
	return *user.Access
}