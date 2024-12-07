package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Add Genre
// @Description Add a new genre to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param data body schemas.GenreAddSchema true "Genre"
// @Success 201 {object} schemas.ResponseSchema "Genre Added Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/genres [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddBookGenres(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.GenreAddSchema{}
	errCode, errData := ValidateRequest(c, &data);
	if errData != nil{
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingGenre := models.Genre{}
	db.Where("LOWER(name) = LOWER(?)", name).First(&existingGenre)
	if existingGenre.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Genre already exists"))
	}
	tags := []models.Tag{}
	if len(tags) > 0 {
		db.Where("slug IN ?", data.TagSlugs).Find(&tags)
		if len(tags) < 1 {
			return c.Status(422).JSON(utils.ValidationErr("tag_slugs", "Enter at leat one valid tag slugs"))
		}
	}
	db.Omit("Tags.*").Create(&models.Genre{Name: name, Tags: tags})
	return c.Status(201).JSON(ResponseMessage("Genre added successfully"))
}

// @Summary Add Tags
// @Description Add a new tag to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param data body schemas.TagsAddSchema true "Tag"
// @Success 201 {object} schemas.ResponseSchema "Tag added successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/tags [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddBookTags(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data);
	if errData != nil{
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingTag := models.Tag{}
	db.Where("LOWER(name) = LOWER(?)", name).First(&existingTag)
	if existingTag.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Tag already exists"))
	}
	db.Create(&models.Tag{Name: name})
	return c.Status(201).JSON(ResponseMessage("Tag added successfully"))
}

// @Summary Update Genre
// @Description Update a genre.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Genre slug"
// @Param data body schemas.GenreAddSchema true "Genre"
// @Success 200 {object} schemas.ResponseSchema "Genre Updated Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/genres/{slug} [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateBookGenre(c *fiber.Ctx) error {
	db := ep.DB
	genre := genreManager.GetBySlug(db, c.Params("slug"))
	if genre == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Genre does not exist"))
	}

	data := schemas.GenreAddSchema{}
	errCode, errData := ValidateRequest(c, &data);
	if errData != nil{
		return c.Status(*errCode).JSON(errData)
	}

	name := data.Name
	existingGenre := models.Genre{}
	db.Where("LOWER(name) = LOWER(?)", name).Not("id = ?", genre.ID).First(&existingGenre)
	if existingGenre.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Genre already exists with that name"))
	}
	tags := []models.Tag{}
	if len(tags) > 0 {
		db.Where("slug IN ?", data.TagSlugs).Find(&tags)
		if len(tags) < 1 {
			return c.Status(422).JSON(utils.ValidationErr("tag_slugs", "Enter at leat one valid tag slugs"))
		}
		genre.Tags = tags
	}
	genre.Name = name
	db.Save(&genre)
	return c.Status(200).JSON(ResponseMessage("Genre updated successfully"))
}

// @Summary Update Tags
// @Description Update a tag to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Tag slug"
// @Param data body schemas.TagsAddSchema true "Tag"
// @Success 200 {object} schemas.ResponseSchema "Tag updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/tags/{slug} [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateBookTags(c *fiber.Ctx) error {
	db := ep.DB

	tag := tagManager.GetBySlug(db, c.Params("slug"))
	if tag == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Tag does not exist"))
	}

	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data);
	if errData != nil{
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingTag := models.Tag{}
	db.Where("LOWER(name) = LOWER(?)", name).Not("id = ?", tag.ID).First(&existingTag)
	if existingTag.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Tag already exists"))
	}
	tag.Name = name
	db.Save(&tag)
	return c.Status(200).JSON(ResponseMessage("Tag updated successfully"))
}

// @Summary List Books with Pagination
// @Description Retrieves a list of books with support for pagination and optional filtering based on book title.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param page query int false "Current Page" default(1)
// @Param title query string false "Title of the book to filter by"
// @Success 200 {object} schemas.BooksResponseSchema "Successfully retrieved list of books"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetBooks(c *fiber.Ctx) error {
	db := ep.DB
	titleQuery := c.Query("title", "")
	books, _ := bookManager.GetLatest(db, "", "", titleQuery)

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

// @Summary List Book Contracts with Pagination
// @Description Retrieves a list of book contracts with support for pagination and optional filtering based on contract status.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param page query int false "Current Page" default(1)
// @Param name query string false "Name of the author to filter by"
// @Param contract_status query string false "status of the contract to filter by" Enums(PENDING, APPROVED, DECLINED, UPDATED)
// @Success 200 {object} schemas.ContractsResponseSchema "Successfully retrieved list of book contracts"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/contracts [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetBookContracts(c *fiber.Ctx) error {
	db := ep.DB
	nameQuery := c.Query("name")
	var name *string
	if nameQuery != "" {
		name = &nameQuery
	}
	contractStatus := choices.ContractStatusChoice(c.Query("contract_status", ""))
	if !contractStatus.IsValid() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid contract status"))
	}

	books := bookManager.GetBookContracts(db, name, &contractStatus)

	// Paginate and return book contracts
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	books = paginatedBooks.([]models.Book)
	response := schemas.ContractsResponseSchema{
		ResponseSchema: ResponseMessage("Book Contracts fetched successfully"),
		Data: schemas.ContractsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(books),
	}
	return c.Status(200).JSON(response)
}
