package routes

import (
	"fmt"

	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	userManager       = managers.UserManager{}
	bookManager       = managers.BookManager{}
	chapterManager    = managers.ChapterManager{}
	boughtBookManager = managers.BoughtBookManager{}
	tagManager        = managers.TagManager{}
	genreManager      = managers.GenreManager{}
	reviewManager     = managers.ReviewManager{}
	replyManager      = managers.ReplyManager{}
	voteManager      = managers.VoteManager{}
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

// @Summary View Latest Books
// @Description This endpoint views a latest books
// @Tags Books
// @Param page query int false "Current Page" default(1)
// @Param genre_slug query string false "Filter by Genre slug"
// @Param tag_slug query string false "Filter by Tag slug"
// @Success 200 {object} schemas.PartialBooksResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books [get]
func (ep Endpoint) GetLatestBooks(c *fiber.Ctx) error {
	db := ep.DB
	genreSlug := c.Query("genre_slug")
	tagSlug := c.Query("tag_slug")
	books, err := bookManager.GetLatest(db, genreSlug, tagSlug, "")
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Paginate and return books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	books = paginatedBooks.([]models.Book)
	response := schemas.PartialBooksResponseSchema{
		ResponseSchema: ResponseMessage("Books fetched successfully"),
		Data: schemas.PartialBooksResponseDataSchema{
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
// @Param tag_slug query string false "Filter by Tag slug"
// @Success 200 {object} schemas.PartialBooksResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/author/{username} [get]
func (ep Endpoint) GetLatestAuthorBooks(c *fiber.Ctx) error {
	db := ep.DB
	username := c.Params("username")
	genreSlug := c.Query("genre_slug")
	tagSlug := c.Query("tag_slug")
	books, err := bookManager.GetLatest(db, genreSlug, tagSlug, "", username)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Paginate and return books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	books = paginatedBooks.([]models.Book)
	response := schemas.PartialBooksResponseSchema{
		ResponseSchema: ResponseMessage("Books fetched successfully"),
		Data: schemas.PartialBooksResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(books),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Single Book (partial chapters)
// @Description This endpoint views a single book
// @Tags Books
// @Param page query int false "Current Page (for reviews pagination)" default(1)
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.PartialBookDetailResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug} [get]
func (ep Endpoint) GetSingleBookPartial(c *fiber.Ctx) error {
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

	book = ViewBook(c, db, *book)

	reviews := paginatedReviews.([]models.Review)
	response := schemas.PartialBookDetailResponseSchema{
		ResponseSchema: ResponseMessage("Book details fetched successfully"),
		Data:           schemas.PartialBookDetailSchema{}.Init(*book, *paginatedData, reviews),
	}
	return c.Status(200).JSON(response)
}

// @Summary View Single Book (All chapters - For The Author and User who has bought it)
// @Description This endpoint views a single book with all chapters
// @Tags Books
// @Param page query int false "Current Page (for reviews pagination)" default(1)
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.BookDetailResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book-full/{slug} [get]
// @Security BearerAuth
func (ep Endpoint) GetSingleBookFull(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	book, err := bookManager.GetBySlugWithReviews(db, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}

	if user.ID != book.AuthorID {
		// Check if current user has bought the book
		boughtBook := boughtBookManager.GetByBuyerAndBook(db, user, *book)
		if boughtBook == nil {
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Only the author or user who has bought the book can see it in full"))
		}
	}

	// Paginate book reviews
	paginatedData, paginatedReviews, err := PaginateQueryset(book.Reviews, c, 30)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	book = ViewBook(c, db, *book)

	reviews := paginatedReviews.([]models.Review)
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
// @Param chapter.title formData string false "First chapter title"
// @Param chapter.text formData string false "First chapter title"
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
		data := map[string]string{
			"genre_slug": "Invalid genre slug!",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Validate Tags
	tagSlugs := data.TagSlugs
	tags, errStr := CheckTagStrings(db, tagSlugs)
	if errStr != nil {
		data := map[string]string{
			"tag_slugs": *errStr,
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Check and validate image
	file, err := ValidateImage(c, "cover_image", true)
	if err != nil {
		return c.Status(422).JSON(err)
	}

	book := bookManager.Create(db, *author, data, genre, "", tags)
	// Upload File
	UploadFile(file, book.CoverImage, cfg.BookCoverImagesBucket)
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
// @Param book formData schemas.BookUpdateSchema true "Book object"
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

	data := schemas.BookUpdateSchema{}
	if errCode, errData := ValidateFormRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate Genre
	genreSlug := data.GenreSlug
	genre := models.Genre{Slug: genreSlug}
	db.Take(&genre, genre)
	if genre.ID == uuid.Nil {
		data := map[string]string{
			"genre_slug": "Invalid genre slug!",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Validate Tags
	tagSlugs := data.TagSlugs
	tags, errStr := CheckTagStrings(db, tagSlugs)
	if errStr != nil {
		data := map[string]string{
			"tag_slugs": *errStr,
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Check and validate image
	file, err := ValidateImage(c, "cover_image", false)
	if err != nil {
		return c.Status(422).JSON(err)
	}

	updatedBook := bookManager.Update(db, *book, data, genre, tags)
	// Upload File
	if file != nil {
		UploadFile(file, updatedBook.CoverImage, cfg.BookCoverImagesBucket)
	}
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
	return c.Status(200).JSON(ResponseMessage("Book deleted successfuly"))
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
	response := schemas.ChapterResponseSchema{
		ResponseSchema: ResponseMessage("Chapter added successfully"),
		Data:           schemas.ChapterSchema{}.Init(chapter),
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
		Data:           schemas.ChapterSchema{}.Init(updatedChapter),
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
	return c.Status(200).JSON(ResponseMessage("Chapter deleted successfuly"))
}

// @Summary Buy A Book
// @Description This endpoint allows a user to buy a book
// @Tags Books
// @Param slug path string true "Book slug"
// @Success 201 {object} schemas.BookResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/book/{slug}/buy [get]
// @Security BearerAuth
func (ep Endpoint) BuyBook(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	book, err := bookManager.GetBySlug(db, c.Params("slug"))
	if err != nil {
		return c.Status(404).JSON(err)
	}

	if user.ID == book.AuthorID {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "You can't buy your own book"))
	}

	bookAlreadyBought := boughtBookManager.GetByBuyerAndBook(db, user, *book)
	if bookAlreadyBought != nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_ALREADY_BOUGHT, "You have bought this book already"))
	}

	if book.Price > user.Coins {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INSUFFICIENT_COINS, "You have insufficient coins"))
	}

	// Create bought book
	boughtBook := boughtBookManager.Create(db, user, *book)

	// Create and send notification in socket
	notification := notificationManager.Create(
		db, user, book.Author, choices.NT_BOOK_PURCHASE,
		fmt.Sprintf("%s bought one of your books.", user.FullName()),
		book, nil, nil, nil,
	)
	SendNotificationInSocket(c, notification)

	response := schemas.BookResponseSchema{
		ResponseSchema: ResponseMessage("Book bought successfully"),
		Data:           schemas.BookSchema{}.Init(boughtBook.Book),
	}
	return c.Status(201).JSON(response)
}

// @Summary View Bought Books
// @Description This endpoint views all books bought by a user
// @Tags Books
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.BooksResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /books/bought [get]
// @Security BearerAuth
func (ep Endpoint) GetBoughtBooks(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	books := boughtBookManager.GetLatest(db, user)
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
// @Router /books/book-full/{slug} [post]
// @Security BearerAuth
func (ep Endpoint) ReviewBook(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	slug := c.Params("slug")
	book, err := bookManager.GetBySlug(db, slug)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Check if current user has bought the book
	boughtBook := boughtBookManager.GetByBuyerAndBook(db, user, *book)
	if boughtBook == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Only the reader who has bought the book can review it"))
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
	text := fmt.Sprintf("%s reviewed your book", user.FullName())
	notification := notificationManager.Create(db, user, boughtBook.Book.Author, choices.NT_REVIEW, text, &boughtBook.Book, &createdReview.ID, nil, nil)
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
// @Router /books/book-full/review/{id} [put]
// @Security BearerAuth
func (ep Endpoint) EditBookReview(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}

	review := reviewManager.GetByUserAndID(db, user, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "You don't have a review with that ID"))
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
// @Router /books/book-full/review/{id} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteBookReview(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}

	review := reviewManager.GetByUserAndID(db, user, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "You don't have a review with that ID"))
	}
	db.Delete(&review)
	return c.Status(200).JSON(ResponseMessage("Review deleted successfully"))
}

// @Summary Get Review Replies
// @Description `This endpoint returns replies of a book review.`
// @Tags Books
// @Param id path string true "Review id (uuid)"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.RepliesResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book-full/review/{id}/replies [get]
func (ep Endpoint) GetReviewReplies(c *fiber.Ctx) error {
	db := ep.DB
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}

	review := reviewManager.GetByID(db, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No review with that ID"))
	}

	// Paginate and return replies
	paginatedData, paginatedReplies, err := PaginateQueryset(review.Replies, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	replies := paginatedReplies.([]models.Reply)
	response := schemas.RepliesResponseSchema{
		ResponseSchema: ResponseMessage("Replies fetched successfully"),
		Data: schemas.RepliesResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(replies),
	}
	return c.Status(200).JSON(response)
}

// @Summary Reply A Review
// @Description `This endpoint allows a user to reply a book review.`
// @Tags Books
// @Param id path string true "Review id (uuid)"
// @Param review body schemas.ReplyReviewSchema true "Reply object"
// @Success 201 {object} schemas.ReplyResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book-full/review/{id}/replies [post]
// @Security BearerAuth
func (ep Endpoint) ReplyReview(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}

	review := reviewManager.GetByID(db, *parsedID)
	if review == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No review with that ID"))
	}

	data := schemas.ReplyReviewSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	reply := replyManager.Create(db, user, review, data)

	// Create and Send Notification in socket
	if user.ID != review.User.ID {
		text := fmt.Sprintf("%s replied your review", user.FullName())
		notification := notificationManager.Create(db, user, review.User, choices.NT_REPLY, text, &review.Book, &review.ID, &reply.ID, nil)
		SendNotificationInSocket(c, notification)
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
// @Param review body schemas.ReplyReviewSchema true "Reply object"
// @Success 200 {object} schemas.ReplyResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /books/book-full/review/replies/{id} [put]
// @Security BearerAuth
func (ep Endpoint) EditReply(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}

	reply := replyManager.GetByUserAndID(db, user, *parsedID)
	if reply == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "You don't have a reply with that ID"))
	}

	data := schemas.ReplyReviewSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	updatedReply := replyManager.Update(db, *reply, data)
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
// @Router /books/book-full/review/replies/{id} [delete]
// @Security BearerAuth
func (ep Endpoint) DeleteReply(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	reviewID := c.Params("id")
	parsedID := ParseUUID(reviewID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}

	reply := replyManager.GetByUserAndID(db, user, *parsedID)
	if reply == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "You don't have a reply with that ID"))
	}
	db.Delete(&reply)
	return c.Status(200).JSON(ResponseMessage("Reply deleted successfully"))
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
	book, err := bookManager.GetBySlug(db, slug)
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
		text := fmt.Sprintf("%s voted your book", user.FullName())
		notification := notificationManager.Create(db, user, book.Author, choices.NT_VOTE, text, book, nil, nil, nil)
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
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid amount parameter"))
	}
	if amount < 1 {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Amount must not be less than 1"))
	}
	if amount > user.Coins {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INSUFFICIENT_COINS, "You have insufficient coins for that conversion"))
	}

	user.Lanterns += amount
	user.Coins -= amount
	db.Save(&user)
	return c.Status(200).JSON(ResponseMessage("Lanterns added successfully"))
}