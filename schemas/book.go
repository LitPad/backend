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

type SectionSchema struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (s SectionSchema) Init(section models.Section) SectionSchema {
	s.Name = section.Name
	s.Slug = section.Slug
	return s
}

type SubSectionSchema struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	BooksCount int    `json:"books_count"`
}

func (s SubSectionSchema) Init(subSection models.SubSection) SubSectionSchema {
	s.Name = subSection.Name
	s.Slug = subSection.Slug
	s.BooksCount = len(subSection.Books)
	return s
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

type ChapterListSchema struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func (c ChapterListSchema) Init(chapter models.Chapter) ChapterListSchema {
	c.Title = chapter.Title
	c.Slug = chapter.Slug
	return c
}

type ParagraphSchema struct {
	Index         uint   `json:"index"`
	Text          string `json:"text"`
	CommentsCount int    `json:"comments_count"`
}

type ChapterDetailSchema struct {
	ChapterListSchema
	Paragraphs []ParagraphSchema `json:"paragraphs"`
}

func (c ChapterDetailSchema) Init(chapter models.Chapter) ChapterDetailSchema {
	c.Title = chapter.Title
	c.Slug = chapter.Slug
	paragraphs := []ParagraphSchema{}
	for _, p := range chapter.Paragraphs {
		paragraphs = append(paragraphs, ParagraphSchema{Text: p.Text, CommentsCount: p.CommentsCount(), Index: p.Index})
	}
	c.Paragraphs = paragraphs
	return c
}

type BookSchema struct {
	Author             UserDataSchema        `json:"author"`
	Title              string                `json:"title"`
	Slug               string                `json:"slug"`
	Blurb              string                `json:"blurb"`
	AgeDiscretion      choices.AgeType       `json:"age_discretion"`
	Genre              GenreWithoutTagSchema `json:"genre"`
	Section            SectionSchema         `json:"section"`
	SubSection         SubSectionSchema      `json:"sub_section"`
	Tags               []TagSchema           `json:"tags"`
	ChaptersCount      int                   `json:"chapters_count"`
	PartialViewChapter *ChapterListSchema    `json:"partial_view_chapter"`
	CoverImage         string                `json:"cover_image"`
	FullPrice          *int                  `json:"full_price"`
	ChapterPrice       int                   `json:"chapter_price"`
	Completed          bool                  `json:"completed"`
	Votes              int                   `json:"votes"`
	Reads              int                   `json:"reads"`
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
	b.Section = b.Section.Init(book.SubSection.Section)
	b.SubSection = b.SubSection.Init(book.SubSection)
	b.ChaptersCount = book.ChaptersCount()
	b.Votes = book.VotesCount()
	b.Reads = book.ReadsCount()
	b.AvgRating = book.AvgRating

	chapters := book.Chapters
	if len(chapters) > 0 {
		chapter := ChapterListSchema{}.Init(chapters[0])
		b.PartialViewChapter = &chapter
	}

	b.CoverImage = book.CoverImage
	b.Completed = book.Completed
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

func (r ReviewSchema) Init(review models.Comment) ReviewSchema {
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
	Text string `json:"text" validate:"required,max=10000"`
}

type CommentSchema struct {
	ParagraphCommentAddSchema
	ID           uuid.UUID      `json:"id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`
	User         UserDataSchema `json:"user"`
	LikesCount   int            `json:"likes_count"`
	RepliesCount int            `json:"replies_count"`
	CreatedAt    time.Time      `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (p CommentSchema) Init(paragraphComment models.Comment) CommentSchema {
	p.ID = paragraphComment.ID
	p.User = p.User.Init(paragraphComment.User)
	p.Text = paragraphComment.Text
	p.LikesCount = paragraphComment.LikesCount()
	p.RepliesCount = paragraphComment.RepliesCount()
	p.CreatedAt = paragraphComment.CreatedAt
	p.UpdatedAt = paragraphComment.UpdatedAt
	return p
}

type ParagraphCommentResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []CommentSchema `json:"items"`
}

type ParagraphCommentResponseSchema struct {
	ResponseSchema
	Data CommentSchema `json:"data"`
}

type ParagraphCommentsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []CommentSchema `json:"comments"`
}

func (p ParagraphCommentsResponseDataSchema) Init(comments []models.Comment) ParagraphCommentsResponseDataSchema {
	// Set Initial Data
	commentItems := make([]CommentSchema, 0)
	for _, comment := range comments {
		commentItems = append(commentItems, CommentSchema{}.Init(comment))
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

func (r RepliesResponseDataSchema) Init(replies []models.Comment) RepliesResponseDataSchema {
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
	WordCount int                       `json:"word_count"`
	Reviews   ReviewsResponseDataSchema `json:"reviews"`
}

func (b BookDetailSchema) Init(book models.Book, reviewsPaginatedData PaginatedResponseDataSchema, reviews []models.Comment) BookDetailSchema {
	b.BookSchema = b.BookSchema.Init(book)
	b.WordCount = book.GetWordCount()
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
	Blurb         string          `form:"blurb" validate:"required,wordcount_max=1000,wordcount_min=100"`
	GenreSlug     string          `form:"genre_slug" validate:"required"`
	TagSlugs      []string        `form:"tag_slugs" validate:"required"`
	AgeDiscretion choices.AgeType `form:"age_discretion" validate:"required,age_discretion_validator"`
}

type ChapterCreateSchema struct {
	Title      string   `json:"title" validate:"required,max=100"`
	Paragraphs []string `json:"paragraphs" validate:"required"`
	IsLast     bool     `json:"is_last"`
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

type SectionsResponseSchema struct {
	ResponseSchema
	Data []SectionSchema `json:"data"`
}

func (s SectionsResponseSchema) Init(sections []models.Section) SectionsResponseSchema {
	// Set Initial Data
	sectionItems := s.Data
	for _, section := range sections {
		sectionItems = append(sectionItems, SectionSchema{}.Init(section))
	}
	s.Data = sectionItems
	return s
}

type SubSectionsResponseSchema struct {
	ResponseSchema
	Data []SubSectionSchema `json:"data"`
}

func (s SubSectionsResponseSchema) Init(subSections []models.SubSection) SubSectionsResponseSchema {
	// Set Initial Data
	subSectionItems := s.Data
	for _, subSection := range subSections {
		subSectionItems = append(subSectionItems, SubSectionSchema{}.Init(subSection))
	}
	s.Data = subSectionItems
	return s
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
	PostalCode           string                       `json:"postal_code"`
	TelephoneNumber      string                       `json:"telephone_number"`
	IDType               choices.ContractIDTypeChoice `json:"id_type"`
	IDFrontImage         string                       `json:"id_front_image"`
	IDBackImage          string                       `json:"id_back_image"`
	BookAvailabilityLink *string                      `json:"book_availability_link"`
	PlannedLength        uint                         `json:"planned_length"`
	AverageChapter       uint                         `json:"average_chapter"`
	UpdateRate           uint                         `json:"update_rate"`
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
	PostalCode           string                       `form:"postal_code" validate:"required"`
	TelephoneNumber      string                       `form:"telephone_number" validate:"required,max=20"`
	IDType               choices.ContractIDTypeChoice `form:"id_type" validate:"required,contract_id_type_validator"`
	BookAvailabilityLink *string                      `form:"book_availability_link"`
	PlannedLength        uint                         `form:"planned_length" validate:"required"`
	AverageChapter       uint                         `form:"average_chapter" validate:"required"`
	UpdateRate           uint                         `form:"update_rate" validate:"required"`
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
	Items []ChapterListSchema `json:"chapters"`
}

func (c ChaptersResponseDataSchema) Init(chapters []models.Chapter) ChaptersResponseDataSchema {
	// Set Initial Data
	chapterItems := make([]ChapterListSchema, 0)
	for _, chapter := range chapters {
		chapterItems = append(chapterItems, ChapterListSchema{}.Init(chapter))
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
	Data ChapterDetailSchema `json:"data"`
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

func (r ReplySchema) Init(reply models.Comment) ReplySchema {
	r.ID = reply.ID
	r.User = r.User.Init(reply.User)
	r.Text = reply.Text
	r.LikesCount = reply.LikesCount()
	r.CreatedAt = reply.CreatedAt
	r.UpdatedAt = reply.UpdatedAt
	return r
}

type BookReportSchema struct {
	Reason                string  `json:"reason" validate:"required,max=1000"`
	AdditionalExplanation *string `json:"additional_explanation" validate:"omitempty,max=1000"`
}
