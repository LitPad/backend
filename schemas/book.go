package schemas

import (
	"github.com/LitPad/backend/models"
	"github.com/google/uuid"
)


type BookSchema struct {
	AuthorID   uuid.UUID      `json:"author_id"`
	Author models.User `json:"author"`
	Title string `json:"title"`
	Blurb  string `json:"blurb"`
	GenreID uuid.UUID `json:"genre_id"`
	Genre models.Genre `json:"genre"`
	Tags []models.Tag `json:"tags"`
	Chapters int `json:"chapters"`
	PartialViewChapters int `json:"partial_view_chapters"`
	WordCount int `json:"word_count"`
	CoverImage string `json:"cover_image"`
	Price int `json:"price"`
	PartialViewFile string `json:"partial_view_file"`
	FullViewFile string `json:"full_view_file"`
}

func (dto *BookSchema) FromModel(book models.Book){
	dto.Author = book.Author
	dto.AuthorID = book.AuthorID
	dto.Blurb = book.Blurb
	dto.Price = book.Price
	dto.Tags = book.Tags
	dto.Title = book.Title
	dto.Genre = book.Genre
	dto.GenreID = book.GenreID
	dto.WordCount = book.WordCount
	dto.Chapters = book.Chapters
	dto.CoverImage = book.CoverImage
	dto.FullViewFile = book.FullViewFile
	dto.PartialViewChapters = book.PartialViewChapters
	dto.PartialViewFile = book.PartialViewFile
}

func(b BookSchema) Init(book models.Book) BookSchema{
	b = BookSchema{
	Author: book.Author,
	AuthorID: book.AuthorID,
	Blurb: book.Blurb,
	Price: book.Price,
	Tags: book.Tags,
	Title: book.Title,
	Genre: book.Genre,
	GenreID: book.GenreID,
	WordCount: book.WordCount,
	Chapters: book.Chapters,
	CoverImage: book.CoverImage,
	FullViewFile: book.FullViewFile,
	PartialViewChapters: book.PartialViewChapters,
	PartialViewFile: book.PartialViewFile,
	}

	return b
}