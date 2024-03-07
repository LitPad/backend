package initials

import (
	"log"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"gorm.io/gorm"
)

func createSuperUser(db *gorm.DB, cfg config.Config) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Admin",
		Email:           cfg.FirstSuperuserEmail,
		Password:        cfg.FirstSuperUserPassword,
		IsSuperuser:     true,
		IsStaff:         true,
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func createWriter(db *gorm.DB, cfg config.Config) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Writer",
		Email:           cfg.FirstWriterEmail,
		Password:        cfg.FirstWriterPassword,
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})

	return user
}

func createReader(db *gorm.DB, cfg config.Config) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Reader",
		Email:           cfg.FirstReaderEmail,
		Password:        cfg.FirstReaderPassword,
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})

	return user
}

func CreateInitialData(db *gorm.DB, cfg config.Config) {
	log.Println("Creating Initial Data....")
	createSuperUser(db, cfg)
	createReader(db, cfg)
	createWriter(db, cfg)
	log.Println("Initial Data Created....")
}