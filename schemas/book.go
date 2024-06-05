package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

type TagSchema struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (t TagSchema) Init(tag models.Tag) TagSchema {
	t.Name = tag.Name
	t.Slug = tag.Slug
	return t
}

type GenreWithoutTagSchema struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (g GenreWithoutTagSchema) Init(genre models.Genre) GenreWithoutTagSchema {
	g.Name = genre.Name
	g.Slug = genre.Slug
	return g
}

type GenreSchema struct {
	GenreWithoutTagSchema
	Tags []TagSchema `json:"tags"`
}

func (g GenreSchema) Init(genre models.Genre) GenreSchema {
	g.GenreWithoutTagSchema = g.GenreWithoutTagSchema.Init(genre)
	tags := genre.Tags
	tagsToAdd := g.Tags
	for _, tag := range tags {
		tagsToAdd = append(tagsToAdd, TagSchema{}.Init(tag))
	}
	g.Tags = tagsToAdd
	return g
}

type ChapterSchema struct {
	Title         string                `json:"title"`
	Slug          string                `json:"slug"`
	Text          string                `json:"text"`
	ChapterStatus choices.ChapterStatus `json:"chapter_status" example:"PUBLISHED"`
	WordCount     int                   `json:"word_count"`
}

func (c ChapterSchema) Init(chapter models.Chapter) ChapterSchema {
	c.Title = chapter.Title
	c.Slug = chapter.Slug
	c.Text = chapter.Text
	c.ChapterStatus = chapter.ChapterStatus
	c.WordCount = chapter.WordCount()
	return c
}

type PartialBookSchema struct {
	Author             UserDataSchema        `json:"author"`
	Title              string                `json:"title"`
	Slug               string                `json:"slug"`
	Blurb              string                `json:"blurb"`
	AgeDiscretion      choices.AgeType       `json:"age_discretion"`
	Genre              GenreWithoutTagSchema `json:"genre"`
	Tags               []TagSchema           `json:"tags"`
	ChaptersCount      int                   `json:"chapters_count"`
	PartialViewChapter *ChapterSchema        `json:"partial_view_chapter"`
	WordCount          int                   `json:"word_count"`
	CoverImage         string                `json:"cover_image"`
	Price              int                   `json:"price"`
	CreatedAt          time.Time             `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	UpdatedAt          time.Time             `json:"updated_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (b PartialBookSchema) Init(book models.Book) PartialBookSchema {
	b.Author = b.Author.Init(book.Author)
	b.Blurb = book.Blurb
	b.Price = book.Price
	b.AgeDiscretion = book.AgeDiscretion

	tags := book.Tags
	tagsToAdd := b.Tags
	for _, tag := range tags {
		tagsToAdd = append(tagsToAdd, TagSchema{}.Init(tag))
	}
	b.Tags = tagsToAdd

	b.Title = book.Title
	b.Slug = book.Slug
	b.Genre = b.Genre.Init(book.Genre)
	b.WordCount = book.WordCount()
	b.ChaptersCount = book.ChaptersCount()

	chapters := book.Chapters
	if len(chapters) > 0 {
		chapter := ChapterSchema{}.Init(chapters[0])
		b.PartialViewChapter = &chapter
	}

	b.CoverImage = book.CoverImage
	b.CreatedAt = book.CreatedAt
	b.UpdatedAt = book.UpdatedAt
	return b
}

type BookSchema struct {
	PartialBookSchema
	Chapters []ChapterSchema `json:"chapters"`
}

func (b BookSchema) Init(book models.Book) BookSchema {
	b.PartialBookSchema = b.PartialBookSchema.Init(book)
	chaptersToAdd := b.Chapters
	chapters := book.Chapters
	for _, chapter := range chapters {
		chaptersToAdd = append(chaptersToAdd, ChapterSchema{}.Init(chapter))
	}
	b.Chapters = chaptersToAdd
	return b
}

type BookChapterCreateSchema struct {
	Title string `json:"title" validate:"required,max=200"`
	Text  string `json:"text" validate:"required,max=100000"`
}

type BookUpdateSchema struct {
	Title         string          `form:"title" validate:"required,max=200"`
	Blurb         string          `form:"blurb" validate:"required,max=200"`
	GenreSlug     string          `form:"genre_slug" validate:"required"`
	TagSlugs      []string        `form:"tag_slugs" validate:"required"`
	Price         int             `form:"price" validate:"required"`
	AgeDiscretion choices.AgeType `form:"age_discretion" validate:"required,age_discretion_validator"`
}

type BookCreateSchema struct {
	BookUpdateSchema
	Chapter *BookChapterCreateSchema `form:"chapter"`
}

type ChapterCreateSchema struct {
	Title         string                `json:"title" validate:"required,max=100"`
	Text          string                `json:"text" validate:"required,max=10000"`
	ChapterStatus choices.ChapterStatus `json:"chapter_status" validate:"required,chapter_status_validator"`
}

type TagsResponseSchema struct {
	ResponseSchema
	Data []TagSchema `json:"data"`
}

func (t TagsResponseSchema) Init(tags []models.Tag) TagsResponseSchema {
	// Set Initial Data
	tagItems := t.Data
	for _, tag := range tags {
		tagItems = append(tagItems, TagSchema{}.Init(tag))
	}
	t.Data = tagItems
	return t
}

type GenresResponseSchema struct {
	ResponseSchema
	Data []GenreSchema `json:"data"`
}

func (g GenresResponseSchema) Init(genres []models.Genre) GenresResponseSchema {
	// Set Initial Data
	genreItems := g.Data
	for _, genre := range genres {
		genreItems = append(genreItems, GenreSchema{}.Init(genre))
	}
	g.Data = genreItems
	return g
}

// Partial Book Responses

type PartialBooksResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []PartialBookSchema `json:"books"`
}

func (b PartialBooksResponseDataSchema) Init(books []models.Book) PartialBooksResponseDataSchema {
	// Set Initial Data
	bookItems := b.Items
	for _, book := range books {
		bookItems = append(bookItems, PartialBookSchema{}.Init(book))
	}
	b.Items = bookItems
	return b
}

type PartialBooksResponseSchema struct {
	ResponseSchema
	Data PartialBooksResponseDataSchema `json:"data"`
}

type PartialBookResponseSchema struct {
	ResponseSchema
	Data PartialBookSchema `json:"data"`
}

// Full Book Responses
type BooksResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []BookSchema `json:"books"`
}

func (b BooksResponseDataSchema) Init(books []models.Book) BooksResponseDataSchema {
	// Set Initial Data
	bookItems := b.Items
	for _, book := range books {
		bookItems = append(bookItems, BookSchema{}.Init(book))
	}
	b.Items = bookItems
	return b
}

type BooksResponseSchema struct {
	ResponseSchema
	Data BooksResponseDataSchema `json:"data"`
}

type BookResponseSchema struct {
	ResponseSchema
	Data BookSchema `json:"data"`
}

type ChapterResponseSchema struct {
	ResponseSchema
	Data ChapterSchema `json:"data"`
}
