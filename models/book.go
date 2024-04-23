package models

import (
	"github.com/google/uuid"
)

type Tag struct {
	BaseModel
	Name string `json:"name"`
}

type Genre struct {
	BaseModel
	Name string `json:"name"`
	Tags []Tag  `json:"tags" gorm:"many2many:genre_tags;"`
}

type Book struct {
	BaseModel
	AuthorID uuid.UUID `json:"author_id"`
	Author   User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Title    string    `json:"title" gorm:"type: varchar(255)"`
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
}

// Note:
// Tags in book must be part of the selected genre
// User can only see allowed chapters, but then the full book will be returned to the user if bought
