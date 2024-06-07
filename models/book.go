package models

import (
	"strings"

	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Tag struct {
	BaseModel
	Name string `gorm:"unique"`
	Slug string `gorm:"unique"`
}

func (tag *Tag) BeforeSave(tx *gorm.DB) (err error) {
	tag.Slug = slug.Make(tag.Name)
	return
}

type Genre struct {
	BaseModel
	Name string `gorm:"unique"`
	Slug string `gorm:"unique"`
	Tags []Tag  `json:"tags" gorm:"many2many:genre_tags;"`
}

func (genre *Genre) BeforeSave(tx *gorm.DB) (err error) {
	genre.Slug = slug.Make(genre.Name)
	return
}

type Book struct {
	BaseModel
	AuthorID      uuid.UUID
	Author        User   `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;<-:false"`
	Title         string `gorm:"type: varchar(255)"`
	Slug          string `gorm:"unique"`
	Blurb         string `gorm:"type: varchar(255)"`
	AgeDiscretion choices.AgeType

	GenreID    uuid.UUID `json:"genre_id"`
	Genre      Genre     `gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE;<-:false"`
	Tags       []Tag     `gorm:"many2many:book_tags;<-:false"`
	Chapters   []Chapter `gorm:"<-:false"`
	CoverImage string    `gorm:"type:varchar(10000)"`

	Price     int      `gorm:"default:0"` // Book price in coins
	Completed bool     `gorm:"default:false"`
	Views     string   `gorm:"type:varchar(10000000)"`
	Reviews   []Review `gorm:"<-:false"`
}

func (b Book) ViewsCount() int {
	views := b.Views
	if len(views) > 0 {
		addresses := strings.Split(b.Views, ", ")
		return len(addresses)
	}
	return 0
}

func (b Book) WordCount() int {
	wordCount := 0
	for _, chapter := range b.Chapters {
		wordCount += chapter.WordCount()
	}
	return wordCount
}

func (b Book) ChaptersCount() int {
	return len(b.Chapters)
}

func (b *Book) GenerateUniqueSlug(tx *gorm.DB) string {
	uniqueSlug := slug.Make(b.Title)
	slug := b.Slug
	if slug != "" {
		uniqueSlug = slug
	}

	existingBook := Book{Slug: uniqueSlug}
	tx.Take(&existingBook, existingBook)
	if existingBook.ID != uuid.Nil && existingBook.ID != b.ID { // slug is already taken
		// Make it unique by attaching a random string
		// to it and repeat the function
		randomStr := utils.GetRandomString(6)
		b.Slug = uniqueSlug + "-" + randomStr
		return b.GenerateUniqueSlug(tx)
	}
	return uniqueSlug
}

func (b *Book) BeforeSave(tx *gorm.DB) (err error) {
	b.Slug = b.GenerateUniqueSlug(tx)
	return
}

type Chapter struct {
	BaseModel
	BookID        uuid.UUID             `json:"book_id"`
	Book          Book                  `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	Title         string                `json:"title" gorm:"type: varchar(255)"`
	Slug          string                `gorm:"unique"`
	Text          string                `json:"text" gorm:"type: varchar(100000)"`
	ChapterStatus choices.ChapterStatus `gorm:"type:varchar(100); default:DRAFT" json:"chapter_status"`
}

func (c *Chapter) GenerateUniqueSlug(tx *gorm.DB) string {
	uniqueSlug := slug.Make(c.Title)
	slug := c.Slug
	if slug != "" {
		uniqueSlug = slug
	}

	existingChapter := Chapter{Slug: uniqueSlug}
	tx.Take(&existingChapter, existingChapter)
	if existingChapter.ID != uuid.Nil && existingChapter.ID != c.ID { // slug is already taken
		// Make it unique by attaching a random string
		// to it and repeat the function
		randomStr := utils.GetRandomString(6)
		c.Slug = uniqueSlug + "-" + randomStr
		return c.GenerateUniqueSlug(tx)
	}
	return uniqueSlug
}

func (c *Chapter) BeforeSave(tx *gorm.DB) (err error) {
	c.Slug = c.GenerateUniqueSlug(tx)
	return
}

func (c Chapter) WordCount() int {
	words := strings.Fields(c.Text)
	wordCount := len(words)
	return wordCount
}

type BoughtBook struct {
	BaseModel
	BuyerID uuid.UUID `gorm:"index:,unique,composite:buyer_id_book_id_bought_book"`
	Buyer   User      `gorm:"foreignKey:BuyerID;constraint:OnDelete:CASCADE;<-:false"`

	BookID uuid.UUID `gorm:"index:,unique,composite:buyer_id_book_id_bought_book"`
	Book   Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
}

type Review struct {
	BaseModel
	UserID uuid.UUID `gorm:"index:,unique,composite:user_id_book_id_reviews"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	BookID uuid.UUID `gorm:"index:,unique,composite:user_id_book_id_reviews"`
	Book   Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`

	Rating  choices.RatingChoice
	Likes   []User `gorm:"many2many:review_likes;<-:false"`
	Text    string
	Replies []Reply `gorm:"<-:false"`
}

func (r Review) LikesCount() int {
	return len(r.Likes)
}

func (r Review) RepliesCount() int {
	return len(r.Replies)
}

type Reply struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	ReviewID uuid.UUID
	Review   Review `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE;<-:false"`

	Likes []User `gorm:"many2many:review_reply_likes;<-:false"`
	Text  string
}

func (r Reply) LikesCount() int {
	return len(r.Likes)
}
