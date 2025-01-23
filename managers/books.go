package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookManager struct {
	Model     models.Book
	ModelList []models.Book
}

func (b BookManager) GetLatest(db *gorm.DB, genreSlug string, tagSlug string, title string, byRating bool, username string, nameContains string) ([]models.Book, *utils.ErrorResponse) {
	books := b.ModelList

	query := db.Model(&b.Model)
	if genreSlug != "" {
		genre := models.Genre{Slug: genreSlug}
		db.Take(&genre, genre)
		if genre.ID == uuid.Nil {
			errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Invalid book genre")
			return books, &errData
		}
		query = query.Where(models.Book{GenreID: genre.ID})
	}
	if tagSlug != "" {
		tag := models.Tag{Slug: tagSlug}
		db.Take(&tag, tag)
		if tag.ID == uuid.Nil {
			return books, nil
		}
		query = query.Where("books.id IN (?)", db.Table("book_tags").Select("book_id").Where("tag_id = ?", tag.ID))
	}

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if username != "" {
		author := models.User{Username: username, AccountType: choices.ACCTYPE_AUTHOR}
		db.Take(&author, author)
		if author.ID == uuid.Nil {
			errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Invalid author username")
			return books, &errData
		}
		query = query.Where(models.Book{AuthorID: author.ID})
	}

	if nameContains != "" {
		query = query.Joins("left join users on users.id = books.author_id").
			Where("users.username ILIKE ? OR users.name ILIKE ?", "%"+nameContains+"%", "%"+nameContains+"%")
	}

	query = query.Select("books.*, COALESCE(AVG(reviews.rating), 0) AS avg_rating").
		Joins("left join reviews on reviews.book_id = books.id").
		Group("books.id")

	if byRating {
		query = query.Order("COALESCE(AVG(reviews.rating), 0) DESC")
	} else {
		query = query.Order("books.created_at DESC")
	}
	query.Scopes(scopes.AuthorGenreTagBookPreloadScope).Find(&books)
	return books, nil
}

func (b BookManager) GetBySlug(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug}
	db.Scopes(scopes.AuthorGenreTagBookScope).Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "No book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetBooksOrderedByRatingAndVotes(db *gorm.DB) []schemas.BookWithStats {
	var books []schemas.BookWithStats

	db.Model(&b.Model).
		Select(`
			books.slug, 
			books.title, 
			books.cover_image, 
			users.username AS author_name, 
			COALESCE(AVG(reviews.rating), 0) AS avg_rating, 
			COUNT(votes.id) AS votes_count, 
			genres.name AS genre_name, 
			genres.slug AS genre_slug
		`).
		Joins("LEFT JOIN users ON users.id = books.author_id"). // Adjust `author_id` if necessary
		Joins("LEFT JOIN reviews ON reviews.book_id = books.id").
		Joins("LEFT JOIN votes ON votes.book_id = books.id").
		Joins("LEFT JOIN genres ON genres.id = books.genre_id"). // Adjust `genre_id` if necessary
		Group("books.slug, books.title, books.cover_image, users.username, genres.name, genres.slug").
		Order("avg_rating DESC, votes_count DESC").
		Limit(10).
		Scan(&books)

	return books
}

func (b BookManager) GetBookContracts(db *gorm.DB, name *string, contractStatus *choices.ContractStatusChoice) []models.Book {
	books := []models.Book{}
	q := db.Not("full_name = ?", "")
	if contractStatus != nil {
		q.Where(models.Book{ContractStatus: *contractStatus})
	}
	if name != nil {
		q.Where(models.Book{FullName: *name})
	}
	q.Find(&books)
	return books
}

func (b BookManager) GetContractedBookBySlug(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug, ContractStatus: choices.CTS_APPROVED}
	db.Scopes(scopes.AuthorGenreTagBookScope).Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "No contract approved book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetBySlugWithReviews(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug}
	db.Scopes(scopes.AuthorGenreTagReviewsBookScope).
		Select("books.*, AVG(reviews.rating) as avg_rating").
		Joins("LEFT JOIN reviews ON reviews.book_id = books.id").
		Group("books.id").
		Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "No book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetByAuthorAndSlug(db *gorm.DB, author *models.User, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{AuthorID: author.ID, Slug: slug}
	db.Scopes(scopes.AuthorGenreTagBookScope).Preload("Chapters").Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "Author has no book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) Create(db *gorm.DB, author models.User, data schemas.BookCreateSchema, genre models.Genre, coverImage string, Tags []models.Tag) models.Book {
	book := models.Book{
		AuthorID: author.ID, Author: author, Title: data.Title,
		Blurb: data.Blurb, AgeDiscretion: data.AgeDiscretion,
		GenreID: genre.ID, Genre: genre,
		Tags:       Tags,
		CoverImage: coverImage,
	}
	db.Omit("Tags.*").Create(&book)
	return book
}

func (b BookManager) Update(db *gorm.DB, book models.Book, data schemas.BookCreateSchema, genre models.Genre, coverImage string, Tags []models.Tag) models.Book {
	book.Title = data.Title
	book.Blurb = data.Blurb
	book.AgeDiscretion = data.AgeDiscretion
	book.GenreID = genre.ID
	book.Genre = genre
	book.Tags = Tags
	if coverImage != "" {
		book.CoverImage = coverImage
	}
	db.Omit("Tags.*").Save(&book)
	return book
}

func (b BookManager) SetContract(db *gorm.DB, book models.Book, idFrontImage string, idBackImage string, data schemas.ContractCreateSchema) models.Book {
	book.FullName = data.FullName
	book.Email = data.Email
	book.PenName = data.PenName
	book.Age = data.Age
	book.Country = data.Country
	book.Address = data.Address
	book.City = data.City
	book.State = data.State
	book.PostalCode = data.PostalCode
	book.TelephoneNumber = data.TelephoneNumber
	book.IDType = data.IDType
	book.BookAvailabilityLink = data.BookAvailabilityLink
	book.PlannedLength = data.PlannedLength
	book.AverageChapter = data.AverageChapter
	book.UpdateRate = data.UpdateRate
	book.Synopsis = data.Synopsis
	book.Outline = data.Outline
	book.IntendedContract = data.IntendedContract
	book.FullPurchaseMode = data.FullPurchaseMode
	if idFrontImage != "" {
		book.IDFrontImage = idFrontImage
	}
	if idBackImage != "" {
		book.IDBackImage = idBackImage
	}

	if book.ContractStatus == choices.CTS_DECLINED {
		book.ContractStatus = choices.CTS_UPDATED
	}
	db.Omit("Tags.*").Save(&book)
	return book
}

type ChapterManager struct {
	Model     models.Chapter
	ModelList []models.Chapter
}

func (c ChapterManager) GetBySlug(db *gorm.DB, slug string) (*models.Chapter, *utils.ErrorResponse) {
	chapter := models.Chapter{Slug: slug}
	db.Joins("Book").Take(&chapter, chapter)
	if chapter.ID == uuid.Nil {
		errD := utils.NotFoundErr("No chapter with that slug")
		return nil, &errD
	}
	return &chapter, nil
}

func (c ChapterManager) GetBySlugWithComments(db *gorm.DB, slug string, index int) (*models.Chapter, *utils.ErrorResponse) {
	chapter := models.Chapter{Slug: slug}
	db.Joins("Book").Preload("Comments", "index = ?", index).Take(&chapter, chapter)
	if chapter.ID == uuid.Nil {
		errD := utils.NotFoundErr("No chapter with that slug")
		return nil, &errD
	}
	return &chapter, nil
}

func (c ChapterManager) IsFirstChapter(db *gorm.DB, chapter models.Chapter) bool {
	firstChapter := c.Model
	db.Order("created_at ASC").First(&firstChapter)
	return firstChapter.ID == chapter.ID
}

func (c ChapterManager) Create(db *gorm.DB, book models.Book, data schemas.ChapterCreateSchema) models.Chapter {
	chapter := models.Chapter{
		BookID: book.ID,
		Title:  data.Title,
		Text:   data.Text,
	}
	db.Create(&chapter)
	return chapter
}

func (c ChapterManager) Update(db *gorm.DB, chapter models.Chapter, data schemas.ChapterCreateSchema) models.Chapter {
	chapter.Title = data.Title
	chapter.Text = data.Text
	db.Save(&chapter)
	return chapter
}

type TagManager struct {
	Model     models.Tag
	ModelList []models.Tag
}

func (t TagManager) GetAll(db *gorm.DB) []models.Tag {
	tags := t.ModelList
	db.Find(&tags)
	return tags
}

func (t TagManager) GetBySlug(db *gorm.DB, slug string) *models.Tag {

	tag := models.Tag{Slug: slug}
	db.Take(&tag, tag)

	if tag.ID == uuid.Nil {
		return nil
	}

	return &tag
}

type GenreManager struct {
	Model     models.Genre
	ModelList []models.Genre
}

func (g GenreManager) GetAll(db *gorm.DB) []models.Genre {
	genres := g.ModelList
	db.Preload("Tags").Find(&genres)
	return genres
}

func (g GenreManager) GetBySlug(db *gorm.DB, slug string) *models.Genre {

	genre := models.Genre{Slug: slug}
	db.Take(&genre, genre)

	if genre.ID == uuid.Nil {
		return nil
	}

	return &genre
}

type ReviewManager struct {
	Model     models.Review
	ModelList []models.Review
}

func (r ReviewManager) GetByID(db *gorm.DB, id uuid.UUID) *models.Review {
	review := r.Model
	db.Where("reviews.id = ?", id).Joins("User").Joins("Book").Preload("Replies").Preload("Replies.User").Preload("Replies.Likes").Take(&review, review)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Review {
	review := models.Review{}
	db.Where("user_id = ?", user.ID).Joins("Book").Joins("User").Preload("Replies").Preload("Likes").Take(&review, id)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) GetByUserAndBook(db *gorm.DB, user *models.User, book models.Book) *models.Review {
	review := models.Review{
		UserID: user.ID,
		BookID: book.ID,
	}
	db.Take(&review, review)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) Create(db *gorm.DB, user *models.User, book models.Book, data schemas.ReviewBookSchema) models.Review {
	review := models.Review{
		UserID: user.ID,
		User:   *user,
		BookID: book.ID,
		Book:   book,
		Rating: data.Rating,
		Text:   data.Text,
	}
	db.Create(&review)
	return review
}

func (r ReviewManager) Update(db *gorm.DB, review models.Review, data schemas.ReviewBookSchema) models.Review {
	review.Text = data.Text
	review.Rating = data.Rating
	db.Save(&review)
	return review
}

type ParagraphCommentManager struct {
	Model     models.ParagraphComment
	ModelList []models.ParagraphComment
}

func (p ParagraphCommentManager) GetByID(db *gorm.DB, id uuid.UUID) *models.ParagraphComment {
	paragraphComment := p.Model
	db.Where("paragraph_comments.id = ?", id).Joins("User").Joins("Chapter").Preload("Replies").Preload("Replies.User").Preload("Replies.Likes").Take(&paragraphComment, paragraphComment)
	if paragraphComment.ID == uuid.Nil {
		return nil
	}
	return &paragraphComment
}

func (p ParagraphCommentManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.ParagraphComment {
	paragraphComment := p.Model
	db.Where("user_id = ?", user.ID).Joins("Chapter").Joins("User").Preload("Replies").Preload("Likes").Take(&paragraphComment, id)
	if paragraphComment.ID == uuid.Nil {
		return nil
	}
	return &paragraphComment
}

func (p ParagraphCommentManager) GetByChapterID(db *gorm.DB, chapterId uuid.UUID) []models.ParagraphComment {
	paragraphComments := p.ModelList
	db.Where("chapter_id = ?", chapterId).Find(&paragraphComments)
	return paragraphComments
}

func (p ParagraphCommentManager) Create(db *gorm.DB, user *models.User, chapterId uuid.UUID, data schemas.ParagraphCommentAddSchema) models.ParagraphComment {
	paragraphComment := models.ParagraphComment{
		UserID:    user.ID,
		User:      *user,
		ChapterID: chapterId,
		Index:     data.Index,
		Text:      data.Text,
	}
	db.Create(&paragraphComment)
	return paragraphComment
}

func (p ParagraphCommentManager) Update(db *gorm.DB, paragraphComment models.ParagraphComment, data schemas.ParagraphCommentAddSchema) models.ParagraphComment {
	paragraphComment.Text = data.Text
	paragraphComment.Index = data.Index
	db.Save(&paragraphComment)
	return paragraphComment
}

type ReplyManager struct {
	Model     models.Reply
	ModelList []models.Reply
}

func (r ReplyManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Reply {
	reply := models.Reply{}
	db.Where("user_id = ?", user.ID).Joins("User").Preload("Likes").Take(&reply, id)
	if reply.ID == uuid.Nil {
		return nil
	}
	return &reply
}

func (r ReplyManager) Create(db *gorm.DB, user *models.User, review *models.Review, paragraphComment *models.ParagraphComment, data schemas.ReplyReviewOrCommentSchema) models.Reply {
	reply := models.Reply{
		UserID: user.ID,
		User:   *user,
		Text:   data.Text,
	}
	if review != nil {
		reply.ReviewID = &review.ID
	} else {
		reply.ParagraphCommentID = &paragraphComment.ID
	}
	db.Create(&reply)
	return reply
}

func (r ReplyManager) Update(db *gorm.DB, reply models.Reply, data schemas.ReplyEditSchema) models.Reply {
	reply.Text = data.Text
	db.Save(&reply)
	return reply
}

type VoteManager struct {
	Model     models.Vote
	ModelList []models.Vote
}

func (v VoteManager) GetByUserAndBook(db *gorm.DB, user *models.User, book *models.Book) *models.Vote {
	vote := models.Vote{UserID: user.ID, BookID: book.ID}
	db.Joins("User").Joins("Book").Take(&vote, vote)
	if vote.ID == uuid.Nil {
		return nil
	}
	return &vote
}

func (v VoteManager) Create(db *gorm.DB, user *models.User, book *models.Book) models.Vote {
	vote := models.Vote{UserID: user.ID, User: *user, Book: *book, BookID: book.ID}
	db.Create(&vote)
	return vote
}
