package routes

import (
	"fmt"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary View Available Book Tags
// @Description This endpoint views available book tags
// @Tags Books
// @Success 200 {object} schemas.TagsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/tags [get]
func (ep Endpoint) GetAllBookTags(c *fiber.Ctx) error {
	db := ep.DB
	tags := tagManager.GetAll(db)

	response := schemas.TagsResponseSchema{
		ResponseSchema: ResponseMessage("Tags fetched successfully"),
	}.Init(tags)
	return c.Status(200).JSON(response)
}

// @Summary View Available Book Genres
// @Description This endpoint views available book genres
// @Tags Books
// @Success 200 {object} schemas.GenresResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/genres [get]
func (ep Endpoint) GetAllBookGenres(c *fiber.Ctx) error {
	db := ep.DB
	genres := genreManager.GetAll(db)

	response := schemas.GenresResponseSchema{
		ResponseSchema: ResponseMessage("Genres fetched successfully"),
	}.Init(genres)
	return c.Status(200).JSON(response)
}

// @Summary View Available Book Sub Genres
// @Description This endpoint views available book sub genres
// @Tags Books
// @Success 200 {object} schemas.SubGenresResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/sub-genres [get]
func (ep Endpoint) GetAllBookSubGenres(c *fiber.Ctx) error {
	db := ep.DB
	genres := genreManager.GetAllSubGenres(db)

	response := schemas.SubGenresResponseSchema{
		ResponseSchema: ResponseMessage("Sub Genres fetched successfully"),
	}.Init(genres)
	return c.Status(200).JSON(response)
}

// @Summary View Latest Books
// @Description This endpoint views a latest books
// @Tags Books
// @Param page query int false "Current Page" default(1)
// @Param genre_slug query string false "Filter by Genre slug"
// @Param sub_genre_slug query string false "Filter by Sub Genre slug"
// @Param tag_slug query string false "Filter by Tag slug"
// @Param featured query bool false "Filter by Featured"
// @Param weeklyFeatured query bool false "Filter by Weekly Featured"
// @Param trending query bool false "Filter by Trending"
// @Success 200 {object} schemas.BooksResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books [get]
func (ep Endpoint) GetLatestBooks(c *fiber.Ctx) error {
	db := ep.DB
	genreSlug := c.Query("genre_slug")
	subGenreSlug := c.Query("sub_genre_slug")
	tagSlug := c.Query("tag_slug")
	featured := c.QueryBool("featured")
	weeklyFeatured := c.QueryBool("weekly_featured")
	trending := c.QueryBool("trending")
	books, err := bookManager.GetLatest(db, genreSlug, subGenreSlug, tagSlug, "", false, "", "", featured, weeklyFeatured, trending)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Paginate and return books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	books = paginatedBooks.([]models.Book)
	response := schemas.BooksResponseSchema{
		ResponseSchema: ResponseMessage("Books fetched successfully"),
		Data: schemas.BooksResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(books),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Latest Books By A Particular Author
// @Description This endpoint views a latest books by an author
// @Tags Books
// @Param page query int false "Current Page" default(1)
// @Param username path string true "Filter by Author Username"
// @Param genre_slug query string false "Filter by Genre slug"
// @Param sub_genre_slug query string false "Filter by Sub Genre slug"
// @Param tag_slug query string false "Filter by Tag slug"
// @Param featured query bool false "Filter by Featured"
// @Param weeklyFeatured query bool false "Filter by Weekly Featured"
// @Param trending query bool false "Filter by Trending"
// @Success 200 {object} schemas.BooksResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/author/{username} [get]
func (ep Endpoint) GetLatestAuthorBooks(c *fiber.Ctx) error {
	db := ep.DB
	username := c.Params("username")
	genreSlug := c.Query("genre_slug")
	subGenreSlug := c.Query("sub_genre_slug")
	tagSlug := c.Query("tag_slug")
	featured := c.QueryBool("featured")
	weeklyFeatured := c.QueryBool("weekly_featured")
	trending := c.QueryBool("trending")
	books, err := bookManager.GetLatest(db, genreSlug, subGenreSlug, tagSlug, "", false, username, "", featured, weeklyFeatured, trending)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Paginate and return books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	books = paginatedBooks.([]models.Book)
	response := schemas.BooksResponseSchema{
		ResponseSchema: ResponseMessage("Books fetched successfully"),
		Data: schemas.BooksResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(books),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Book Chapters
// @Description `This endpoint views chapters of a book`
// @Description `A Guest user will view just the first chapter`
// @Description `An Authenticated user will view all the chapters if he's subscribed or he gets only the first chapter`
// @Description `The owner will view all chapters of the book`
// @Tags Books
// @Param slug path string true "Get Chapter by Book Slug"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.ChaptersResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug}/chapters [get]
// @Security BearerAuth
func (ep Endpoint) GetBookChapters(c *fiber.Ctx) error {
	db := ep.DB
	slug := c.Params("slug")
	book, err := bookManager.GetBySlug(db, slug, true)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	paginatedData, paginatedChapters, err := PaginateQueryset(book.Chapters, c, 50)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	chapters := paginatedChapters.([]models.Chapter)
	response := schemas.ChaptersResponseSchema{
		ResponseSchema: ResponseMessage("Chapters fetched successfully"),
		Data: schemas.ChaptersResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(chapters),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Book Chapter
// @Description `This endpoint views a single chapter of a book`
// @Description `An inactive subscriber can only view the chapter if its the first one`
// @Tags Books
// @Param slug path string true "Get Chapter by Slug"
// @Success 200 {object} schemas.ChapterResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/chapters/chapter/{slug} [get]
// @Security BearerAuth
func (ep Endpoint) GetBookChapter(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	chapter, err := chapterManager.GetBySlug(db, slug)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	chapterIsFirst := chapterManager.IsFirstChapter(db, *chapter)
	if chapter.Book.AuthorID != user.ID && user.SubscriptionExpired() && !chapterIsFirst && !user.IsStaff {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Renew your subscription to view this chapter"))
	}
	ReadBook(db, chapter.BookID, user)
	response := schemas.ChapterResponseSchema{
		ResponseSchema: ResponseMessage("Chapter fetched successfully"),
		Data:           schemas.ChapterDetailSchema{}.Init(*chapter),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Comments Of A Paragraph of A Chapter
// @Description `This endpoint view comments of a single paragraph of a chapter`
// @Description `An inactive subscriber can only view the paragraph comment if its the first one`
// @Tags Books
// @Param slug path string true "Chapter Slug"
// @Param index path int true "Paragraph Index"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.ParagraphCommentsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/chapters/chapter/{slug}/paragraph/{index}/comments [get]
// @Security BearerAuth
func (ep Endpoint) GetParagraphComments(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	index, _ := c.ParamsInt("index", 1)
	chapter, comments, err := chapterManager.GetBySlugWithComments(db, slug, uint(index))
	if err != nil {
		return c.Status(404).JSON(err)
	}
	chapterIsFirst := chapterManager.IsFirstChapter(db, *chapter)
	if chapter.Book.AuthorID != user.ID && user.SubscriptionExpired() && !chapterIsFirst && !user.IsStaff {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Renew your subscription to view this chapter"))
	}

	// Paginate and return comments
	paginatedData, paginatedComments, err := PaginateQueryset(comments, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	comments = paginatedComments.([]models.Comment)

	response := schemas.ParagraphCommentsResponseSchema{
		ResponseSchema: ResponseMessage("Paragraph Comments fetched successfully"),
		Data: schemas.ParagraphCommentsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(comments),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Single Book
// @Description This endpoint views a single book
// @Tags Books
// @Param page query int false "Current Page (for reviews pagination)" default(1)
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.BookDetailResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug} [get]
func (ep Endpoint) GetSingleBook(c *fiber.Ctx) error {
	db := ep.DB
	book, err := bookManager.GetBySlugWithReviews(db, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Paginate book reviews
	paginatedData, paginatedReviews, err := PaginateQueryset(book.Reviews, c, 30)
	if err != nil {
		return c.Status(400).JSON(err)
	}

	reviews := paginatedReviews.([]models.Comment)
	response := schemas.BookDetailResponseSchema{
		ResponseSchema: ResponseMessage("Book details fetched successfully"),
		Data:           schemas.BookDetailSchema{}.Init(*book, *paginatedData, reviews),
	}
	return c.Status(200).JSON(response)
}

// @Summary Create A Book
// @Description This endpoint allows a writer to create a book
// @Tags Books
// @Param book formData schemas.BookCreateSchema true "Book object"
// @Param cover_image formData file true "Cover Image to upload"
// @Success 201 {object} schemas.BookResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books [post]
// @Security BearerAuth
func (ep Endpoint) CreateBook(c *fiber.Ctx) error {
	db := ep.DB
	author := RequestUser(c)
	data := schemas.BookCreateSchema{}
	if errCode, errData := ValidateFormRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate Genre
	genreSlug := data.GenreSlug
	genre := models.Genre{Slug: genreSlug}
	db.Take(&genre, genre)
	if genre.ID == uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("genre_slug", "Invalid genre slug!"))
	}

	// Validate Tags
	tagSlugs := data.TagSlugs
	tags, errStr := CheckTagStrings(db, tagSlugs)
	if errStr != nil {
		return c.Status(422).JSON(utils.ValidationErr("tag_slugs", *errStr))
	}

	// Check and validate image
	file, err := ValidateImage(c, "cover_image", true)
	if err != nil {
		return c.Status(422).JSON(err)
	}

	// Upload File
	coverImage := UploadFile(file, string(choices.IF_BOOKS))
	book := bookManager.Create(db, *author, data, genre, coverImage, tags)
	response := schemas.BookResponseSchema{
		ResponseSchema: ResponseMessage("Book created successfully"),
		Data:           schemas.BookSchema{}.Init(book),
	}
	return c.Status(201).JSON(response)
}

// @Summary Update A Book
// @Description This endpoint allows a writer to update a book
// @Tags Books
// @Param slug path string true "Book slug"
// @Param book formData schemas.BookCreateSchema true "Book object"
// @Param cover_image formData file false "Cover Image to upload"
// @Success 200 {object} schemas.BookResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug} [put]
// @Security BearerAuth
func (ep Endpoint) UpdateBook(c *fiber.Ctx) error {
	db := ep.DB
	author := RequestUser(c)
	book, err := bookManager.GetByAuthorAndSlug(db, author, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}

	data := schemas.BookCreateSchema{}
	if errCode, errData := ValidateFormRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate Genre
	genreSlug := data.GenreSlug
	genre := models.Genre{Slug: genreSlug}
	db.Take(&genre, genre)
	if genre.ID == uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("genre_slug", "Invalid genre slug!"))
	}

	// Validate Tags
	tagSlugs := data.TagSlugs
	tags, errStr := CheckTagStrings(db, tagSlugs)
	if errStr != nil {
		return c.Status(422).JSON(utils.ValidationErr("tag_slugs", *errStr))
	}

	// Check and validate image
	file, err := ValidateImage(c, "cover_image", false)
	if err != nil {
		return c.Status(422).JSON(err)
	}

	// Upload File
	coverImage := ""
	if file != nil {
		coverImage = UploadFile(file, string(choices.IF_BOOKS))
	}

	updatedBook := bookManager.Update(db, *book, data, genre, coverImage, tags)

	response := schemas.BookResponseSchema{
		ResponseSchema: ResponseMessage("Book updated successfully"),
		Data:           schemas.BookSchema{}.Init(updatedBook),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete A Book
// @Description This endpoint allows a writer to delete a book
// @Tags Books
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteBook(c *fiber.Ctx) error {
	db := ep.DB
	author := RequestUser(c)
	book, err := bookManager.GetByAuthorAndSlug(db, author, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}
	db.Delete(&book)
	return c.Status(200).JSON(ResponseMessage("Book deleted successfully"))
}

// @Summary Add A Chapter to a Book
// @Description `This endpoint allows a writer to add a chapter to his/her book`
// @Description `Chapter status: DRAFT, PUBLISHED, TRASH`
// @Tags Books
// @Param slug path string true "Book slug"
// @Param chapter body schemas.ChapterCreateSchema true "Chapter object"
// @Success 201 {object} schemas.ChapterResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug}/add-chapter [post]
// @Security BearerAuth
func (ep Endpoint) AddChapter(c *fiber.Ctx) error {
	db := ep.DB
	author := RequestUser(c)
	book, err := bookManager.GetByAuthorAndSlug(db, author, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}

	data := schemas.ChapterCreateSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	chapter := chapterManager.Create(db, *book, data)
	if data.IsLast {
		book.Completed = true
		db.Save(&book)
	}
	response := schemas.ChapterResponseSchema{
		ResponseSchema: ResponseMessage("Chapter added successfully"),
		Data:           schemas.ChapterDetailSchema{}.Init(chapter),
	}
	return c.Status(201).JSON(response)
}

// @Summary Update A Chapter of a Book
// @Description `This endpoint allows a writer to update a chapter in his/her book`
// @Description `Chapter status: DRAFT, PUBLISHED, TRASH`
// @Tags Books
// @Param slug path string true "Chapter slug"
// @Param chapter body schemas.ChapterCreateSchema true "Chapter object"
// @Success 200 {object} schemas.ChapterResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/chapter/{slug} [put]
// @Security BearerAuth
func (ep Endpoint) UpdateChapter(c *fiber.Ctx) error {
	db := ep.DB
	author := RequestUser(c)
	chapter, err := chapterManager.GetBySlug(db, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}

	if chapter.Book.AuthorID != author.ID {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to edit"))
	}

	data := schemas.ChapterCreateSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	updatedChapter := chapterManager.Update(db, *chapter, data)
	response := schemas.ChapterResponseSchema{
		ResponseSchema: ResponseMessage("Chapter updated successfully"),
		Data:           schemas.ChapterDetailSchema{}.Init(updatedChapter),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete A Chapter
// @Description This endpoint allows a writer to delete a chapter from a book
// @Tags Books
// @Param slug path string true "Chapter slug"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/chapter/{slug} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteChapter(c *fiber.Ctx) error {
	db := ep.DB
	author := RequestUser(c)
	chapter, err := chapterManager.GetBySlug(db, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}
	if chapter.Book.AuthorID != author.ID {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to delete"))
	}
	db.Delete(&chapter)
	return c.Status(200).JSON(ResponseMessage("Chapter deleted successfully"))
}

// @Summary Review A Book
// @Description `This endpoint allows a user to review a book.`
// @Description `The author cannot review his own book.`
// @Description `Only the reader who has bought the book can review the book.`
// @Description `A reader cannot add multiple reviews to a book.`
// @Tags Books
// @Param slug path string true "Book slug"
// @Param review body schemas.ReviewBookSchema true "Review object"
// @Success 201 {object} schemas.ReviewResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug} [post]
// @Security BearerAuth
func (ep Endpoint) ReviewBook(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	book, err := bookManager.GetBySlug(db, slug, true)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Check if current user has bought at least a chapter of the book
	if user.SubscriptionExpired() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "User doesn't have active subscription"))
	}
	data := schemas.ReviewBookSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	review := reviewManager.GetByUserAndBook(db, user, *book)
	if review != nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_ALREADY_REVIEWED, "This book has been reviewed by you already"))
	}

	createdReview := reviewManager.Create(db, user, *book, data)

	// Create and Send Notification in socket
	text := fmt.Sprintf("%s reviewed your book", user.Username)
	notification := notificationManager.Create(db, user, book.Author, choices.NT_REVIEW, text, book, &createdReview.ID, nil)
	SendNotificationInSocket(c, notification)

	response := schemas.ReviewResponseSchema{
		ResponseSchema: ResponseMessage("Review created successfully"),
		Data:           schemas.ReviewSchema{}.Init(createdReview),
	}
	return c.Status(201).JSON(response)
}

// @Summary Edit Book Review
// @Description `This endpoint allows a user to edit his/her book review.`
// @Tags Books
// @Param id path string true "Review id (uuid)"
// @Param review body schemas.ReviewBookSchema true "Review object"
// @Success 200 {object} schemas.ReviewResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/review/{id} [put]
// @Security BearerAuth
func (ep Endpoint) EditBookReview(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	review := reviewManager.GetByUserAndID(db, user, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.NotFoundErr("You don't have a review with that ID"))
	}
	data := schemas.ReviewBookSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	updatedReview := reviewManager.Update(db, *review, data)
	response := schemas.ReviewResponseSchema{
		ResponseSchema: ResponseMessage("Review updated successfully"),
		Data:           schemas.ReviewSchema{}.Init(updatedReview),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete Book Review
// @Description `This endpoint allows a user to delete his/her book review.`
// @Tags Books
// @Param id path string true "Review id (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/review/{id} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteBookReview(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	review := reviewManager.GetByUserAndID(db, user, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.NotFoundErr("You don't have a review with that ID"))
	}
	db.Delete(&review)
	return c.Status(200).JSON(ResponseMessage("Review deleted successfully"))
}

// @Summary Get Comment/Review Replies
// @Description `This endpoint returns replies of a book review or paragraph comment`
// @Tags Books
// @Param id path string true "Comment/Review id (uuid)"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.RepliesResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/review/{id}/replies [get]
func (ep Endpoint) GetReviewReplies(c *fiber.Ctx) error {
	db := ep.DB
	commentID := c.Params("id")
	parsedID := ParseUUID(commentID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	review := reviewManager.GetByID(db, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.NotFoundErr("No review or comment with that ID"))
	}

	// Paginate and return replies
	paginatedData, paginatedReplies, err := PaginateQueryset(review.Replies, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	replies := paginatedReplies.([]models.Comment)
	response := schemas.RepliesResponseSchema{
		ResponseSchema: ResponseMessage("Replies fetched successfully"),
		Data: schemas.RepliesResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(replies),
	}
	return c.Status(200).JSON(response)
}

// @Summary Reply A Review Or A Paragraph Comment
// @Description `This endpoint allows a user to reply a book review.`
// @Tags Books
// @Param id path string true "Review or Paragraph Comment id (uuid)"
// @Param review body schemas.ReplyReviewOrCommentSchema true "Reply object"
// @Success 201 {object} schemas.ReplyResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/review-or-paragraph-comment/{id}/replies [post]
// @Security BearerAuth
func (ep Endpoint) ReplyReviewOrParagraphComment(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewOrParagraphCommentID := c.Params("id")
	parsedID := ParseUUID(reviewOrParagraphCommentID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	data := schemas.ReplyReviewOrCommentSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	var reply models.Comment
	if data.Type == choices.RT_REVIEW {
		review := reviewManager.GetByID(db, *parsedID)
		if review == nil {
			return c.Status(404).JSON(utils.NotFoundErr("No review with that ID"))
		}
		reply = commentManager.CreateReply(db, user, review, data)
		// Create and Send Notification in socket
		if user.ID != review.User.ID {
			text := fmt.Sprintf("%s replied your review", user.Username)
			notification := notificationManager.Create(db, user, review.User, choices.NT_REPLY, text, review.Book, &review.ID, nil)
			SendNotificationInSocket(c, notification)
		}
	} else {
		paragraphComment := commentManager.GetByID(db, *parsedID, false)
		if paragraphComment == nil {
			return c.Status(404).JSON(utils.NotFoundErr("No paragraph comment with that ID"))
		}
		reply = commentManager.CreateReply(db, user, paragraphComment, data)
	}

	response := schemas.ReplyResponseSchema{
		ResponseSchema: ResponseMessage("Reply created successfully"),
		Data:           schemas.ReplySchema{}.Init(reply),
	}
	return c.Status(201).JSON(response)
}

// @Summary Edit A Reply
// @Description `This endpoint allows a user to edit his/her reply`
// @Tags Books
// @Param id path string true "Reply id (uuid)"
// @Param review body schemas.ReplyEditSchema true "Reply object"
// @Success 200 {object} schemas.ReplyResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/review-or-paragraph-comment/replies/{id} [put]
// @Security BearerAuth
func (ep Endpoint) EditReply(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	replyID := c.Params("id")
	parsedID := ParseUUID(replyID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	reply := commentManager.GetReplyByUserAndID(db, user, *parsedID)
	if reply == nil {
		return c.Status(404).JSON(utils.NotFoundErr("You don't have a reply with that ID"))
	}

	data := schemas.ReplyEditSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	updatedReply := commentManager.UpdateReply(db, *reply, data)
	response := schemas.ReplyResponseSchema{
		ResponseSchema: ResponseMessage("Reply updated successfully"),
		Data:           schemas.ReplySchema{}.Init(updatedReply),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete A Reply
// @Description `This endpoint allows a user to delete his/her reply`
// @Tags Books
// @Param id path string true "Reply id (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/review-or-paragraph-comment/replies/{id} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteReply(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	reply := commentManager.GetReplyByUserAndID(db, user, *parsedID)
	if reply == nil {
		return c.Status(404).JSON(utils.NotFoundErr("You don't have a reply with that ID"))
	}
	db.Delete(&reply)
	return c.Status(200).JSON(ResponseMessage("Reply deleted successfully"))
}

// @Summary Add A Comment To A Paragraph In A Book Chapter
// @Description `This endpoint allows a user to add a comment in a paragraph to a book chapter.`
// @Tags Books
// @Param slug path string true "Chapter slug"
// @Param index path int true "Paragraph Index of the chapter"
// @Param review body schemas.ParagraphCommentAddSchema true "Paragraph Comment object"
// @Success 201 {object} schemas.ParagraphCommentResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/chapters/chapter/{slug}/paragraph/{index}/comments [post]
// @Security BearerAuth
func (ep Endpoint) AddParagraphComment(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	index, _ := c.ParamsInt("index", 1)
	if index < 1 {
		return c.Status(400).JSON(utils.InvalidParamErr("Enter a valid index"))
	}

	chapter, err := chapterManager.GetBySlug(db, slug)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	paragraph := chapterManager.GetParagraph(db, *chapter, uint(index))
	if paragraph == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Paragraph does not exist"))
	}

	data := schemas.ParagraphCommentAddSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	paragraphComment := commentManager.Create(db, user, paragraph.ID, data)
	response := schemas.ParagraphCommentResponseSchema{
		ResponseSchema: ResponseMessage("Comment created successfully"),
		Data:           schemas.CommentSchema{}.Init(paragraphComment),
	}
	return c.Status(201).JSON(response)
}

// @Summary Edit Paragraph Comment
// @Description `This endpoint allows a user to edit his/her paragraph comment.`
// @Tags Books
// @Param id path string true "Comment id (uuid)"
// @Param review body schemas.ParagraphCommentAddSchema true "Comment object"
// @Success 200 {object} schemas.ParagraphCommentResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/chapters/chapter/paragraph-comment/{id} [put]
// @Security BearerAuth
func (ep Endpoint) EditParagraphComment(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	commentID := c.Params("id")
	parsedID := ParseUUID(commentID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	paragraphComment := commentManager.GetByUserAndID(db, user, *parsedID)
	if paragraphComment == nil {
		return c.Status(404).JSON(utils.NotFoundErr("You don't have a comment with that ID"))
	}
	data := schemas.ParagraphCommentAddSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	updatedComment := commentManager.Update(db, *paragraphComment, data)
	response := schemas.ParagraphCommentResponseSchema{
		ResponseSchema: ResponseMessage("Comment updated successfully"),
		Data:           schemas.CommentSchema{}.Init(updatedComment),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete Paragraph Comment
// @Description `This endpoint allows a user to delete his/her paragraph comment.`
// @Tags Books
// @Param id path string true "Review id (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/chapters/chapter/paragraph-comment/{id} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteParagraphComment(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	commentID := c.Params("id")
	parsedID := ParseUUID(commentID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("You entered an invalid uuid"))
	}

	comment := commentManager.GetByUserAndID(db, user, *parsedID)
	if comment == nil {
		return c.Status(404).JSON(utils.NotFoundErr("You don't have a comment with that ID"))
	}
	db.Delete(&comment)
	return c.Status(200).JSON(ResponseMessage("Comment deleted successfully"))
}

// @Summary Vote A Book
// @Description This endpoint allows a user to vote a book
// @Tags Books
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug}/vote [get]
// @Security BearerAuth
func (ep Endpoint) VoteBook(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	book, err := bookManager.GetBySlug(db, slug, true)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Check if user has enough lanterns to vote
	if user.Lanterns < 1 {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INSUFFICIENT_LANTERNS, "You have insufficient lanterns to vote"))
	}
	createdVote := voteManager.Create(db, user, book)
	// Create and Send Notification in socket
	if user.ID != createdVote.UserID {
		text := fmt.Sprintf("%s voted your book", user.Username)
		notification := notificationManager.Create(db, user, book.Author, choices.NT_VOTE, text, book, nil, nil)
		SendNotificationInSocket(c, notification)
	}
	user.Lanterns -= 1
	db.Save(&user)
	return c.Status(200).JSON(ResponseMessage("Book voted successfully"))
}

// @Summary Convert Coins To Lanterns
// @Description This endpoint allows a user to convert coins to lanterns
// @Tags Books
// @Param amount path int true "Amount to convert"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/lanterns-generation/{amount} [get]
// @Security BearerAuth
func (ep Endpoint) ConvertCoinsToLanterns(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	amount, err := c.ParamsInt("amount")
	if err != nil {
		return c.Status(400).JSON(utils.InvalidParamErr("Invalid amount parameter"))
	}
	if amount < 1 {
		return c.Status(400).JSON(utils.InvalidParamErr("Amount must not be less than 1"))
	}
	if amount > user.Coins {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INSUFFICIENT_COINS, "You have insufficient coins for that conversion"))
	}

	user.Lanterns += amount
	user.Coins -= amount
	db.Save(&user)
	return c.Status(200).JSON(ResponseMessage("Lanterns added successfully"))
}

// @Summary Set Contract
// @Description `This endpoint allows a user to create/update a contract for his/her book`
// @Tags Books
// @Param slug path string true "Book slug"
// @Param contract formData schemas.ContractCreateSchema true "Contract object"
// @Param id_front_image formData file false "Front Image of your id"
// @Param id_back_image formData file false "Back Image of your id"
// @Success 200 {object} schemas.ContractResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book/{slug}/set-contract [post]
// @Security BearerAuth
func (ep Endpoint) SetContract(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	book, err := bookManager.GetByAuthorAndSlug(db, user, slug)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	if book.ContractStatus == choices.CTS_APPROVED {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_CONTRACT_ALREADY_APPROVED, "This book already has an approved contract"))
	}
	data := schemas.ContractCreateSchema{}
	if errCode, errData := ValidateFormRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Check and validate image
	imageRequired := false
	if book.FullName == "" {
		imageRequired = true
	}
	idFrontImageFile, idFrontImageFileErr := ValidateImage(c, "id_front_image", imageRequired)
	if idFrontImageFileErr != nil {
		return c.Status(422).JSON(idFrontImageFileErr)
	}

	idBackImageFile, idBackImageFileErr := ValidateImage(c, "id_back_image", imageRequired)
	if idBackImageFileErr != nil {
		return c.Status(422).JSON(idBackImageFileErr)
	}

	// Upload File
	var idFrontImage string
	var idBackImage string
	if idFrontImageFile != nil {
		idFrontImage = UploadFile(idFrontImageFile, "ID_FRONT_IMAGES")
	}
	if idBackImageFile != nil {
		idBackImage = UploadFile(idBackImageFile, "ID_BACK_IMAGES")
	}

	updatedBook := bookManager.SetContract(db, *book, idFrontImage, idBackImage, data)
	response := schemas.ContractResponseSchema{
		ResponseSchema: ResponseMessage("Contract set successfully"),
		Data:           schemas.ContractSchema{}.Init(updatedBook),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Bookmarked Books
// @Description This endpoint allows a user to view his/her bookmarked books
// @Tags Books
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.BooksResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/bookmarked [get]
// @Security BearerAuth
func (ep Endpoint) GetBookmarkedBooks(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	books := bookManager.GetUserBookmarkedBooks(db, *user)
	// Paginate and return books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	books = paginatedBooks.([]models.Book)
	response := schemas.BooksResponseSchema{
		ResponseSchema: ResponseMessage("Books fetched successfully"),
		Data: schemas.BooksResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(books),
	}
	return c.Status(200).JSON(response)
}

// @Summary Bookmark A Book
// @Description This endpoint allows a user to bookmark a book
// @Tags Books
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug}/bookmark [get]
// @Security BearerAuth
func (ep Endpoint) BookmarkBook(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	book, err := bookManager.GetBySlug(db, slug, false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	status := bookmarkManager.AddOrDelete(db, *user, *book)
	return c.Status(200).JSON(ResponseMessage(status + " successfully"))
}

// @Summary Report A Book
// @Description This endpoint allows a user to report a book
// @Tags Books
// @Param slug path string true "Book slug"
// @Param report body schemas.BookReportSchema true "Report object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug}/report [post]
// @Security BearerAuth
func (ep Endpoint) ReportBook(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	book, err := bookManager.GetBySlug(db, slug, false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	data := schemas.BookReportSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	bookReportManager.Create(db, *user, *book, data.Reason)
	return c.Status(200).JSON(ResponseMessage("Report submitted successfully"))
}

// @Summary Like/Unlike A Comment/Reply
// @Description `This endpoint allows a user to like/unlike a comment or a reply (a kind of toggle)`
// @Tags Books
// @Param id path string true "Comment or reply id (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/chapters/chapter/comment/{id} [get]
// @Security BearerAuth
func (ep Endpoint) LikeAComment(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	id := ParseUUID(c.Params("id"))
	if id == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("Enter a valid uuid"))
	}
	commentOrReply := commentManager.GetByID(db, *id, false)
	if commentOrReply == nil {
		return c.Status(404).JSON(err)
	}
	status := likeManager.AddOrDelete(db, *user, *commentOrReply)
	return c.Status(200).JSON(ResponseMessage(status + " successfully"))
}
