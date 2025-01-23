package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/google/uuid"
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
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Text  string `json:"text"`
	Trash bool   `json:"trash" example:"false"`
}

func (c ChapterSchema) Init(chapter models.Chapter) ChapterSchema {
	c.Title = chapter.Title
	c.Slug = chapter.Slug
	c.Text = chapter.Text
	c.Trash = chapter.Trash
	return c
}

type BookSchema struct {
	Author             UserDataSchema        `json:"author"`
	Title              string                `json:"title"`
	Slug               string                `json:"slug"`
	Blurb              string                `json:"blurb"`
	AgeDiscretion      choices.AgeType       `json:"age_discretion"`
	Genre              GenreWithoutTagSchema `json:"genre"`
	Tags               []TagSchema           `json:"tags"`
	ChaptersCount      int                   `json:"chapters_count"`
	PartialViewChapter *ChapterSchema        `json:"partial_view_chapter"`
	CoverImage         string                `json:"cover_image"`
	FullPrice          *int                  `json:"full_price"`
	ChapterPrice       int                   `json:"chapter_price"`
	Views              int                   `json:"views"`
	Votes              int                   `json:"votes"`
	AvgRating          float64               `json:"avg_rating"`
	CreatedAt          time.Time             `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	UpdatedAt          time.Time             `json:"updated_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (b BookSchema) Init(book models.Book) BookSchema {
	b.Author = b.Author.Init(book.Author)
	b.Blurb = book.Blurb
	b.FullPrice = book.FullPrice
	b.ChapterPrice = book.ChapterPrice
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
	b.ChaptersCount = book.ChaptersCount()
	b.Votes = book.VotesCount()
	b.AvgRating = book.AvgRating

	chapters := book.Chapters
	if len(chapters) > 0 {
		chapter := ChapterSchema{}.Init(chapters[0])
		b.PartialViewChapter = &chapter
	}

	b.CoverImage = book.CoverImage
	b.Views = book.ViewsCount()
	b.CreatedAt = book.CreatedAt
	b.UpdatedAt = book.UpdatedAt
	return b
}

type ReviewBookSchema struct {
	Rating choices.RatingChoice `json:"rating" validate:"required,rating_choice_validator"`
	Text   string               `json:"text" validate:"required,max=10000"`
}

type ReviewSchema struct {
	ReviewBookSchema
	ID           uuid.UUID      `json:"id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`
	User         UserDataSchema `json:"user"`
	LikesCount   int            `json:"likes_count"`
	RepliesCount int            `json:"replies_count"`
	CreatedAt    time.Time      `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (r ReviewSchema) Init(review models.Review) ReviewSchema {
	r.ID = review.ID
	r.User = r.User.Init(review.User)
	r.Rating = review.Rating
	r.Text = review.Text
	r.LikesCount = review.LikesCount()
	r.RepliesCount = review.RepliesCount()
	r.CreatedAt = review.CreatedAt
	r.UpdatedAt = review.UpdatedAt
	return r
}

type ReviewsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []ReviewSchema `json:"items"`
}

type ReviewResponseSchema struct {
	ResponseSchema
	Data ReviewSchema `json:"data"`
}

type ParagraphCommentAddSchema struct {
	Index int    `json:"paragraph_index" validate:"required"`
	Text  string `json:"text" validate:"required,max=10000"`
}

type ParagraphCommentSchema struct {
	ParagraphCommentAddSchema
	ID           uuid.UUID      `json:"id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`
	User         UserDataSchema `json:"user"`
	LikesCount   int            `json:"likes_count"`
	RepliesCount int            `json:"replies_count"`
	CreatedAt    time.Time      `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (p ParagraphCommentSchema) Init(paragraphComment models.ParagraphComment) ParagraphCommentSchema {
	p.ID = paragraphComment.ID
	p.User = p.User.Init(paragraphComment.User)
	p.Index = paragraphComment.Index
	p.Text = paragraphComment.Text
	p.LikesCount = paragraphComment.LikesCount()
	p.RepliesCount = paragraphComment.RepliesCount()
	p.CreatedAt = paragraphComment.CreatedAt
	p.UpdatedAt = paragraphComment.UpdatedAt
	return p
}

type ParagraphCommentResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []ParagraphCommentSchema `json:"items"`
}

type ParagraphCommentResponseSchema struct {
	ResponseSchema
	Data ParagraphCommentSchema `json:"data"`
}

type ParagraphCommentsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []ParagraphCommentSchema `json:"replies"`
}

func (p ParagraphCommentsResponseDataSchema) Init(comments []models.ParagraphComment) ParagraphCommentsResponseDataSchema {
	// Set Initial Data
	commentItems := make([]ParagraphCommentSchema, 0)
	for _, comment := range comments {
		commentItems = append(commentItems, ParagraphCommentSchema{}.Init(comment))
	}
	p.Items = commentItems
	return p
}

type ParagraphCommentsResponseSchema struct {
	ResponseSchema
	Data ParagraphCommentsResponseDataSchema `json:"data"`
}

type RepliesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []ReplySchema `json:"replies"`
}

func (r RepliesResponseDataSchema) Init(replies []models.Reply) RepliesResponseDataSchema {
	// Set Initial Data
	replyItems := make([]ReplySchema, 0)
	for _, reply := range replies {
		replyItems = append(replyItems, ReplySchema{}.Init(reply))
	}
	r.Items = replyItems
	return r
}

type RepliesResponseSchema struct {
	ResponseSchema
	Data RepliesResponseDataSchema `json:"data"`
}

type ReplyResponseSchema struct {
	ResponseSchema
	Data ReplySchema `json:"data"`
}

type BookDetailSchema struct {
	BookSchema
	Reviews ReviewsResponseDataSchema `json:"reviews"`
}

func (b BookDetailSchema) Init(book models.Book, reviewsPaginatedData PaginatedResponseDataSchema, reviews []models.Review) BookDetailSchema {
	b.BookSchema = b.BookSchema.Init(book)
	reviewsToAdd := b.Reviews.Items
	for _, review := range reviews {
		reviewsToAdd = append(reviewsToAdd, ReviewSchema{}.Init(review))
	}
	b.Reviews = ReviewsResponseDataSchema{
		PaginatedResponseDataSchema: reviewsPaginatedData,
		Items:                       reviewsToAdd,
	}
	return b
}

type BookDetailResponseSchema struct {
	ResponseSchema
	Data BookDetailSchema `json:"data"`
}

type BookChapterCreateSchema struct {
	Title string `json:"title" validate:"required,max=200"`
	Text  string `json:"text" validate:"required,max=100000"`
}

type BookCreateSchema struct {
	Title         string          `form:"title" validate:"required,max=200"`
	Blurb         string          `form:"blurb" validate:"required,max=200"`
	GenreSlug     string          `form:"genre_slug" validate:"required"`
	TagSlugs      []string        `form:"tag_slugs" validate:"required"`
	AgeDiscretion choices.AgeType `form:"age_discretion" validate:"required,age_discretion_validator"`
}

type ChapterCreateSchema struct {
	Title string `json:"title" validate:"required,max=100"`
	Text  string `json:"text" validate:"required,max=10000"`
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

type ContractSchema struct {
	FullName             string                       `json:"full_name"`
	Email                string                       `json:"email"`
	PenName              string                       `json:"pen_name"`
	Age                  uint                         `json:"age"`
	Country              string                       `json:"country"`
	Address              string                       `json:"address"`
	City                 string                       `json:"city"`
	State                string                       `json:"state"`
	PostalCode           uint                         `json:"postal_code"`
	TelephoneNumber      string                       `json:"telephone_number"`
	IDType               choices.ContractIDTypeChoice `json:"id_type"`
	IDFrontImage         string                       `json:"id_front_image"`
	IDBackImage          string                       `json:"id_back_image"`
	BookAvailabilityLink *string                      `json:"book_availability_link"`
	PlannedLength        uint                         `json:"planned_length"`
	AverageChapter       uint                         `json:"average_chapter"`
	UpdateRate           uint                         `json:"update_rate"`
	Synopsis             string                       `json:"synopsis"`
	Outline              string                       `json:"outline"`
	IntendedContract     choices.ContractTypeChoice   `json:"intended_contract"`
	FullPurchaseMode     bool                         `json:"full_purchase_mode"`
	ContractStatus       choices.ContractStatusChoice `json:"contract_status"`
	FullPrice            *int                         `json:"full_price"`
	ChapterPrice         int                          `json:"chapter_price"`
}

func (c ContractSchema) Init(book models.Book) ContractSchema {
	c.FullName = book.FullName
	c.Email = book.Email
	c.PenName = book.PenName
	c.Age = book.Age
	c.Country = book.Country
	c.Address = book.Address
	c.City = book.City
	c.State = book.State
	c.PostalCode = book.PostalCode
	c.TelephoneNumber = book.TelephoneNumber
	c.IDType = book.IDType
	c.BookAvailabilityLink = book.BookAvailabilityLink
	c.PlannedLength = book.PlannedLength
	c.AverageChapter = book.AverageChapter
	c.UpdateRate = book.UpdateRate
	c.Synopsis = book.Synopsis
	c.Outline = book.Outline
	c.IntendedContract = book.IntendedContract
	c.FullPurchaseMode = book.FullPurchaseMode
	c.ContractStatus = book.ContractStatus
	c.FullPrice = book.FullPrice
	c.ChapterPrice = book.ChapterPrice
	c.IDFrontImage = book.IDFrontImage
	c.IDBackImage = book.IDBackImage
	return c
}

type ContractCreateSchema struct {
	FullName             string                       `form:"full_name" validate:"required,max=1000"`
	Email                string                       `form:"email" validate:"required,email"`
	PenName              string                       `form:"pen_name" validate:"required,max=1000"`
	Age                  uint                         `form:"age" validate:"required"`
	Country              string                       `form:"country" validate:"required,max=1000"`
	Address              string                       `form:"address" validate:"required,max=1000"`
	City                 string                       `form:"city" validate:"required,max=1000"`
	State                string                       `form:"state" validate:"required,max=1000"`
	PostalCode           uint                         `form:"postal_code" validate:"required"`
	TelephoneNumber      string                       `form:"telephone_number" validate:"required,max=20"`
	IDType               choices.ContractIDTypeChoice `form:"id_type" validate:"required,contract_id_type_validator"`
	BookAvailabilityLink *string                      `form:"book_availability_link"`
	PlannedLength        uint                         `form:"planned_length" validate:"required"`
	AverageChapter       uint                         `form:"average_chapter" validate:"required"`
	UpdateRate           uint                         `form:"update_rate" validate:"required"`
	Synopsis             string                       `form:"synopsis" validate:"required"`
	Outline              string                       `form:"outline" validate:"required"`
	IntendedContract     choices.ContractTypeChoice   `form:"intended_contract" validate:"required,contract_type_validator"`
	FullPurchaseMode     bool                         `form:"full_purchase_mode"`
}

// Book Responses
type BooksResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []BookSchema `json:"books"`
}

func (b BooksResponseDataSchema) Init(books []models.Book) BooksResponseDataSchema {
	// Set Initial Data
	bookItems := make([]BookSchema, 0)
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

type ContractResponseSchema struct {
	ResponseSchema
	Data ContractSchema `json:"data"`
}

type ChaptersResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []ChapterSchema `json:"chapters"`
}

func (c ChaptersResponseDataSchema) Init(chapters []models.Chapter) ChaptersResponseDataSchema {
	// Set Initial Data
	chapterItems := make([]ChapterSchema, 0)
	for _, chapter := range chapters {
		chapterItems = append(chapterItems, ChapterSchema{}.Init(chapter))
	}
	c.Items = chapterItems
	return c
}

type ChaptersResponseSchema struct {
	ResponseSchema
	Data ChaptersResponseDataSchema `json:"data"`
}

type ChapterResponseSchema struct {
	ResponseSchema
	Data ChapterSchema `json:"data"`
}

type ReplyEditSchema struct {
	Text string `json:"text" validate:"required,max=10000"`
}

type ReplyReviewOrCommentSchema struct {
	ReplyEditSchema
	Type choices.ReplyType `json:"type" validate:"required,reply_type_validator"`
}

type ReplySchema struct {
	ID         uuid.UUID      `json:"id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`
	User       UserDataSchema `json:"user"`
	Text       string         `json:"text"`
	LikesCount int            `json:"likes_count"`
	CreatedAt  time.Time      `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	UpdatedAt  time.Time      `json:"updated_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (r ReplySchema) Init(reply models.Reply) ReplySchema {
	r.ID = reply.ID
	r.User = r.User.Init(reply.User)
	r.Text = reply.Text
	r.LikesCount = reply.LikesCount()
	r.CreatedAt = reply.CreatedAt
	r.UpdatedAt = reply.UpdatedAt
	return r
}
