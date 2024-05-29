package schemas

import (
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
	for i := range tags {
		tagsToAdd = append(tagsToAdd, TagSchema{}.Init(tags[i]))
	}
	g.Tags = tagsToAdd
	return g
}

type ChapterSchema struct {
	Title         string                `json:"title"`
	Text          string                `json:"text"`
	ChapterStatus choices.ChapterStatus `json:"chapter_status" example:"PUBLISHED"`
	WordCount     int                   `json:"word_count"`
}

func (c ChapterSchema) Init(chapter models.Chapter) ChapterSchema {
	c.Title = chapter.Title
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
	Genre              GenreWithoutTagSchema `json:"genre"`
	Tags               []TagSchema           `json:"tags"`
	ChaptersCount      int                   `json:"chapters_count"`
	PartialViewChapter *ChapterSchema         `json:"partial_view_chapter"`
	WordCount          int                   `json:"word_count"`
	CoverImage         string                `json:"cover_image"`
	Price              int                   `json:"price"`
}

func (b PartialBookSchema) Init(book models.Book) PartialBookSchema {
	b.Author = b.Author.Init(book.Author)
	b.Blurb = book.Blurb
	b.Price = book.Price

	tags := book.Tags
	tagsToAdd := b.Tags
	for i := range tags {
		tagsToAdd = append(tagsToAdd, TagSchema{}.Init(tags[i]))
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
	return b
}

type BookSchema struct {
	PartialBookSchema
	Chapters []ChapterSchema `json:"chapters"`
	File     string          `json:"file"`
}

func (b BookSchema) Init(book models.Book) BookSchema {
	b.PartialBookSchema = b.PartialBookSchema.Init(book)
	chaptersToAdd := b.Chapters
	chapters := book.Chapters
	for i := range chapters {
		chaptersToAdd = append(chaptersToAdd, ChapterSchema{}.Init(chapters[i]))
	}
	b.Chapters = chaptersToAdd
	b.File = book.File
	return b
}

type BookCreateSchema struct {
	Title               string   `json:"title"`
	Blurb               string   `json:"blurb"`
	GenreSlug           string   `json:"genre_slug"`
	TagSlugs            []string `json:"tag_slugs"`
	Chapters            int      `json:"chapters"`
	PartialViewChapters int      `json:"partial_view_chapters"`
	WordCount           int      `json:"word_count"`
	CoverImage          string   `json:"cover_image"`
	Price               int      `json:"price"`
	PartialViewFile     string   `json:"partial_view_file"`
	FullViewFile        string   `json:"full_view_file"`
}

type TagsResponseSchema struct {
	ResponseSchema
	Data []TagSchema `json:"data"`
}

func (t TagsResponseSchema) Init(tags []models.Tag) TagsResponseSchema {
	// Set Initial Data
	tagItems := t.Data
	for i := range tags {
		tagItems = append(tagItems, TagSchema{}.Init(tags[i]))
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
	for i := range genres {
		genreItems = append(genreItems, GenreSchema{}.Init(genres[i]))
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
	for i := range books {
		bookItems = append(bookItems, PartialBookSchema{}.Init(books[i]))
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
	for i := range books {
		bookItems = append(bookItems, BookSchema{}.Init(books[i]))
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
