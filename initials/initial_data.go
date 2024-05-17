package initials

import (
	"math/rand"
	"log"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func createSuperUser(db *gorm.DB, cfg config.Config) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Admin",
		Username:        "test-admin",
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
		Username:        "test-writer",
		AccountType:     choices.ACCTYPE_WRITER,
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
		Username:        "test-reader",
		Email:           cfg.FirstReaderEmail,
		Password:        cfg.FirstReaderPassword,
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})

	return user
}

func createCoins(db *gorm.DB) {
	coins := []models.Coin{}
	db.Find(&coins)
	if len(coins) < 1 {
		for i := 1; i <= 10; i++ {
			defaultPrice := decimal.NewFromFloat(20.25)
			coin := models.Coin{Amount: 10 * i, Price: defaultPrice.Mul(decimal.NewFromInt(int64(i)))}
			coins = append(coins, coin)
		}
		db.Create(&coins)
	}
}

func createTags(db *gorm.DB) []models.Tag {
	tags := []models.Tag{}
	db.Find(&tags)
	if len(tags) < 1 {
		for i := range TAGS {
			tag := models.Tag{Name: TAGS[i]}
			tags = append(tags, tag)
		}
		db.Create(&tags)
	}
	return tags
}

func createGenres(db *gorm.DB, tags []models.Tag) {
	genres := []models.Genre{}
	db.Find(&genres)

	if len(genres) < 1 {
		for i := range GENRES {
			// Shuffle the list
			rand.New(rand.NewSource(time.Now().UnixNano()))
			rand.Shuffle(len(tags), func(i, j int) {
				tags[i], tags[j] = tags[j], tags[i]
			})
			genre := models.Genre{Name: GENRES[i], Tags: tags[:10]}
			genres = append(genres, genre)
		}
		db.Omit("Tags.*").Create(&genres)
	}
}

func CreateInitialData(db *gorm.DB, cfg config.Config) {
	log.Println("Creating Initial Data....")
	createSuperUser(db, cfg)
	createReader(db, cfg)
	createWriter(db, cfg)
	createCoins(db)
	tags := createTags(db)
	createGenres(db, tags)
	log.Println("Initial Data Created....")
}
