package models

import (
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
	AuthorID uuid.UUID `json:"author_id"`
	Author   User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Title    string    `json:"title" gorm:"type: varchar(255)"`
	Slug     string    `gorm:"unique"`
	Blurb    string    `json:"blurb" gorm:"type: varchar(255)"`

	GenreID             uuid.UUID `json:"genre_id"`
	Genre               Genre     `gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE"`
	Tags                []Tag     `json:"tags" gorm:"many2many:book_tags;"`
	Chapters            int       `gorm:"default:0"`
	PartialViewChapters int       `gorm:"default:0"` // Amount of chapters allowed to view freely
	WordCount           int       `gorm:"default:0" json:"word_count"`
	CoverImage          string    `gorm:"type:varchar(10000)" json:"cover_image"`

	Price           int    `gorm:"default:0"`           // Book price in coins
	PartialViewFile string `gorm:"type:varchar(10000)"` // Partial File to view
	FullViewFile    string `gorm:"type:varchar(10000)"` // Full File to view
	Completed       bool   `gorm:"default:false"`
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

// Note:
// Tags in book must be part of the selected genre
// User can only see allowed chapters, but then the full book will be returned to the user if bought
