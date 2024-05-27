package schemas

import (
	"github.com/LitPad/backend/models"
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

type BookSchema struct {
	Author              UserDataSchema        `json:"author"`
	Title               string                `json:"title"`
	Slug                string                `json:"slug"`
	Blurb               string                `json:"blurb"`
	Genre               GenreWithoutTagSchema `json:"genre"`
	Tags                []TagSchema           `json:"tags"`
	Chapters            int                   `json:"chapters"`
	PartialViewChapters int                   `json:"partial_view_chapters"`
	WordCount           int                   `json:"word_count"`
	CoverImage          string                `json:"cover_image"`
	Price               int                   `json:"price"`
	PartialViewFile     string                `json:"partial_view_file"`
	FullViewFile        string                `json:"full_view_file"`
}

func (b BookSchema) Init(book models.Book) BookSchema {
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
	b.WordCount = book.WordCount
	b.Chapters = book.Chapters
	b.CoverImage = book.CoverImage
	bookFile := book.FullViewFile
	b.FullViewFile = bookFile
	if bookFile == "" {
		b.FullViewFile = "Payment required"
	}
	b.PartialViewChapters = book.PartialViewChapters
	b.PartialViewFile = book.PartialViewFile
	return b
}

type BookCreateSchema struct {
	Title               string      `json:"title"`
	Blurb               string      `json:"blurb"`
	GenreSlug           string      `json:"genre_slug"`
	TagSlugs            []string `json:"tag_slugs"`
	Chapters            int         `json:"chapters"`
	PartialViewChapters int         `json:"partial_view_chapters"`
	WordCount           int         `json:"word_count"`
	CoverImage          string      `json:"cover_image"`
	Price               int         `json:"price"`
	PartialViewFile     string      `json:"partial_view_file"`
	FullViewFile        string      `json:"full_view_file"`
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
