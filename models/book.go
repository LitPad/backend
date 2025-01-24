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
	Name   string  `gorm:"unique"`
	Slug   string  `gorm:"unique"`
	Genres []Genre `gorm:"many2many:genre_tags;"`
}

func (tag *Tag) BeforeSave(tx *gorm.DB) (err error) {
	tag.Slug = slug.Make(tag.Name)
	return
}

type Genre struct {
	BaseModel
	Name string `gorm:"unique"`
	Slug string `gorm:"unique"`
	Tags []Tag  `gorm:"many2many:genre_tags;"`
}

func (genre *Genre) BeforeSave(tx *gorm.DB) (err error) {
	genre.Slug = slug.Make(genre.Name)
	return
}

type Book struct {
	BaseModel
	AuthorID      uuid.UUID
	Author        User   `gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL;<-:false"`
	Title         string `gorm:"type: varchar(255)"`
	Slug          string `gorm:"unique"`
	Blurb         string `gorm:"type: varchar(255)"`
	AgeDiscretion choices.AgeType

	GenreID    uuid.UUID `json:"genre_id"`
	Genre      Genre     `gorm:"foreignKey:GenreID;constraint:OnDelete:SET NULL;<-:false"`
	Tags       []Tag     `gorm:"many2many:book_tags;<-:false"`
	Chapters   []Chapter `gorm:"<-:false"`
	CoverImage string    `gorm:"type:varchar(10000)"`

	Completed bool     `gorm:"default:false"`
	Views     string   `gorm:"type:varchar(10000000)"`
	Reviews   []Review `gorm:"<-:false"`
	Votes     []Vote   `gorm:"<-:false"`

	AvgRating float64 // meant for query purposes. do not intentionally populate field

	// BOOK CONTRACT
	FullName             string `gorm:"type: varchar(1000)"`
	Email                string
	PenName              string `gorm:"type: varchar(1000)"`
	Age                  uint
	Country              string `gorm:"type: varchar(1000)"`
	Address              string `gorm:"type: varchar(1000)"`
	City                 string `gorm:"type: varchar(1000)"`
	State                string `gorm:"type: varchar(1000)"`
	PostalCode           string
	TelephoneNumber      string                       `gorm:"type: varchar(20)"`
	IDType               choices.ContractIDTypeChoice `gorm:"type: varchar(100)"`
	IDFrontImage         string
	IDBackImage          string
	BookAvailabilityLink *string
	PlannedLength        uint
	AverageChapter       uint
	UpdateRate           uint
	Synopsis             string
	Outline              string
	IntendedContract     choices.ContractTypeChoice
	FullPrice            *int
	ChapterPrice         int
	FullPurchaseMode     bool                         `gorm:"default:false"`
	ContractStatus       choices.ContractStatusChoice `gorm:"default:PENDING"`
}

func (b Book) ViewsCount() int {
	views := b.Views
	if len(views) > 0 {
		addresses := strings.Split(b.Views, ", ")
		return len(addresses)
	}
	return 0
}

func (b Book) VotesCount() int {
	return len(b.Votes)
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

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	slug := b.GenerateUniqueSlug(tx)
	b.Slug = slug
	return
}

type Chapter struct {
	BaseModel
	BookID   uuid.UUID          `json:"book_id"`
	Book     Book               `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	Title    string             `json:"title" gorm:"type: varchar(255)"`
	Slug     string             `gorm:"unique"`
	Text     string             `gorm:"type:text"`
	Trash    bool               `gorm:"default:false"`
	Comments []ParagraphComment `gorm:"<-:false"`
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

type ParagraphComment struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	ChapterID uuid.UUID
	Chapter   Chapter `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE;<-:false"`
	Index     int
	Likes     []User  `gorm:"many2many:paragraph_comment_likes;<-:false"`
	Text      string  `gorm:"type:varchar(10000)"`
	Replies   []Reply `gorm:"<-:false"`
}

func (p ParagraphComment) LikesCount() int {
	return len(p.Likes)
}

func (p ParagraphComment) RepliesCount() int {
	return len(p.Replies)
}

// REVIEWS
type Review struct {
	BaseModel
	UserID uuid.UUID `gorm:"index:,unique,composite:user_id_book_id_reviews"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	BookID uuid.UUID `gorm:"index:,unique,composite:user_id_book_id_reviews"`
	Book   Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`

	Rating  choices.RatingChoice
	Likes   []User  `gorm:"many2many:review_likes;<-:false"`
	Text    string  `gorm:"type:varchar(10000)"`
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

	ReviewID *uuid.UUID
	Review   *Review `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE;<-:false"`

	ParagraphCommentID *uuid.UUID
	ParagraphComment   *ParagraphComment `gorm:"foreignKey:ParagraphCommentID;constraint:OnDelete:CASCADE;<-:false"`

	Likes []User `gorm:"many2many:review_or_paragraph_comment_reply_likes;<-:false"`
	Text  string `gorm:"type:varchar(10000)"`
}

func (r Reply) LikesCount() int {
	return len(r.Likes)
}

type Vote struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	BookID uuid.UUID
	Book   Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
}
