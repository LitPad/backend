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

func createAuthor(db *gorm.DB, cfg config.Config) models.User {
	name := "Dark Xenia"
	user := models.User{
		Name:            &name,
		Username:        "dark-xenia",
		AccountType:     choices.ACCTYPE_AUTHOR,
		Email:           cfg.FirstAuthorEmail,
		Password:        cfg.FirstAuthorPassword,
		IsEmailVerified: true,
		IsSuperuser:     true,
		IsStaff:         true,
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

func createGenres(db *gorm.DB, tags []models.Tag) []models.Genre {
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
	return genres
}

func createSubGenres(db *gorm.DB) []models.SubGenre {
	subGenres := []models.SubGenre{}
	db.Find(&subGenres)

	if len(subGenres) < 1 {
		for _, item := range SUBGENRES {
			subGenre := models.SubGenre{Name: item}
			subGenres = append(subGenres, subGenre)
		}
		db.Create(&subGenres)
	}
	return subGenres
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

func createBook(db *gorm.DB, author models.User, genre models.Genre, tag models.Tag, subGenre models.SubGenre) models.Book {
	book := models.Book{}
	bookToCreate := BOOK
	bookToCreate.AuthorID = author.ID
	bookToCreate.GenreID = genre.ID
	bookToCreate.SubGenreID = subGenre.ID
	bookToCreate.Tags = []models.Tag{tag}
	db.Omit("Tags.*").FirstOrCreate(&book, bookToCreate)
	return book
}

func createChapter(db *gorm.DB, book models.Book) models.Chapter {
	chapter := models.Chapter{
		BookID: book.ID, Title: "The Journey Begins",
	}
	db.Preload("Paragraphs").FirstOrCreate(&chapter, chapter)
	return chapter
}

func createParagraphs(db *gorm.DB, chapter models.Chapter) models.Paragraph {
	paragraphs := chapter.Paragraphs
	if len(paragraphs) < 1 {
		paragraphsToCreate := []models.Paragraph{}
		for idx, paragraph := range PARAGRAPHS {
			paragraphsToCreate = append(paragraphsToCreate, models.Paragraph{ChapterID: chapter.ID, Index: uint(idx + 1), Text: paragraph})
		}
		db.Create(&paragraphsToCreate)
		paragraphs = paragraphsToCreate
	}
	return paragraphs[0]
}

func createParagraphComment(db *gorm.DB, user models.User, paragraph models.Paragraph) models.Comment {
	comment := models.Comment{
		UserID: user.ID, ParagraphID: &paragraph.ID,
		Text: "Wow, he's in trouble",
	}
	db.FirstOrCreate(&comment, comment)
	return comment
}

func createReply(db *gorm.DB, comment models.Comment, user models.User) models.Comment {
	reply := models.Comment{
		UserID: user.ID, ParentID: &comment.ID,
		Text: "Wow, you're right",
	}
	db.FirstOrCreate(&reply, reply)
	return reply
}

func createReview(db *gorm.DB, book models.Book, user models.User) models.Comment {
	review := models.Comment{
		UserID: user.ID, BookID: &book.ID, Rating: choices.RC_5,
		Text: "This is the best book I've ever read.",
	}
	db.FirstOrCreate(&review, review)
	return review
}

func CreateInitialData(db *gorm.DB, cfg config.Config) {
	log.Println("Creating Initial Data....")
	createSuperUser(db, cfg)
	createReader(db, cfg)
	author := createAuthor(db, cfg)
	createCoins(db)
	tags := createTags(db)
	genres := createGenres(db, tags)
	subGenres := createSubGenres(db)
	createGifts(db)
	createSubscriptionPlans(db)
	book := createBook(db, author, genres[0], tags[0], subGenres[0])
	chapter := createChapter(db, book)
	paragraph := createParagraphs(db, chapter)
	comment := createParagraphComment(db, author, paragraph)
	createReply(db, comment, author)
	createReview(db, book, author)
	log.Println("Initial Data Created....")
}
