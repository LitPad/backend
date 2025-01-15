package initials

import (
	"log"
	"math/rand"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func createSuperUser(db *gorm.DB, cfg config.Config) models.User {
	name := "Test Admin"
	user := models.User{
		Name:            &name,
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
	name := "Test Author"
	user := models.User{
		Name:            &name,
		Username:        "test-author",
		AccountType:     choices.ACCTYPE_AUTHOR,
		Email:           cfg.FirstAuthorEmail,
		Password:        cfg.FirstAuthorPassword,
		IsEmailVerified: true,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})

	return user
}

func createReader(db *gorm.DB, cfg config.Config) models.User {
	name := "Test Reader"
	user := models.User{
		Name:            &name,
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
		for _, item := range TAGS {
			tag := models.Tag{Name: item}
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
		for _, item := range GENRES {
			// Shuffle the list
			rand.New(rand.NewSource(time.Now().UnixNano()))
			rand.Shuffle(len(tags), func(i, j int) {
				tags[i], tags[j] = tags[j], tags[i]
			})
			genre := models.Genre{Name: item, Tags: tags[:10]}
			genres = append(genres, genre)
		}
		db.Omit("Tags.*").Create(&genres)
	}
}

func createGifts(db *gorm.DB) {
	gifts := []models.Gift{}
	db.Find(&gifts)
	if len(gifts) < 1 {
		for i, name := range GIFTNAMES {
			gift := models.Gift{Name: name, Price: 100 * i, Image: "https://img.com", Lanterns: 2 * i}
			gifts = append(gifts, gift)
		}
		db.Create(&gifts)
	}
}

func createSubscriptionPlans(db *gorm.DB) {
	plans := []models.SubscriptionPlan{}
	db.Find(&plans)
	if len(plans) < 1 {
		monthlyAmount, _ := decimal.NewFromString("12.99")
		annualAmount, _ := decimal.NewFromString("131.88")
		plansToCreate := []*models.SubscriptionPlan{
			{Amount: monthlyAmount, SubType: choices.ST_MONTHLY},
			{Amount: annualAmount, SubType: choices.ST_ANNUAL},
		}
		db.Create(&plansToCreate)
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
	createGifts(db)
	createSubscriptionPlans(db)
	log.Println("Initial Data Created....")
}
