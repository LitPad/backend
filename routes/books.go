package routes

import (
	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	bookManager  = managers.BookManager{}
	chapterManager  = managers.ChapterManager{}
	tagManager   = managers.TagManager{}
	genreManager = managers.GenreManager{}
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
	books, err := bookManager.GetLatest(db, genreSlug, tagSlug)
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
	books, err := bookManager.GetLatest(db, genreSlug, tagSlug, username)
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
	fileUrl, err := ValidateAndUploadImage(c, "cover_image", "books", true)
	if err != nil {
		return c.Status(422).JSON(err)
	}

	book := bookManager.Create(db, *author, data, genre, *fileUrl, tags)
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
	fileUrl, err := ValidateAndUploadImage(c, "cover_image", "books", false)
	if err != nil {
		return c.Status(422).JSON(err)
	}

	updatedBook := bookManager.Update(db, *book, data, genre, fileUrl, tags)
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