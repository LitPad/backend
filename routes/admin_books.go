package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Summary Add Genre
// @Description Add a new genre to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param data body schemas.TagsAddSchema true "Genre"
// @Success 201 {object} schemas.ResponseSchema "Genre Added Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/genres [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddBookGenre(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingGenre := models.Genre{}
	db.Where("LOWER(name) = LOWER(?)", name).First(&existingGenre)
	if existingGenre.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Genre already exists"))
	}
	db.Create(&models.Genre{Name: name})
	return c.Status(201).JSON(ResponseMessage("Genre added successfully"))
}

// @Summary Add Section
// @Description Add a new book section to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param data body schemas.TagsAddSchema true "Section"
// @Success 201 {object} schemas.ResponseSchema "Section Added Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/sections [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddBookSection(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingSection := models.Section{}
	db.Where("LOWER(name) = LOWER(?)", name).First(&existingSection)
	if existingSection.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Section already exists"))
	}
	db.Create(&models.Section{Name: name})
	return c.Status(201).JSON(ResponseMessage("Section added successfully"))
}

// @Summary Add SubSection
// @Description Add a new book subsection to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Section slug"
// @Param data body schemas.TagsAddSchema true "SubSection"
// @Success 201 {object} schemas.ResponseSchema "Sub Section Added Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/sections/{slug}/subsections [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddBookSubSection(c *fiber.Ctx) error {
	db := ep.DB
	section := genreManager.GetSectionBySlug(db, c.Params("slug"))
	if section == nil {
		return c.Status(404).JSON(utils.NotFoundErr("No section with that id"))
	}
	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingSubSection := models.SubSection{}
	db.Where("LOWER(name) = LOWER(?)", name).First(&existingSubSection)
	if existingSubSection.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Sub Section already exists"))
	}
	db.Create(&models.SubSection{Name: name, SectionID: section.ID})
	return c.Status(201).JSON(ResponseMessage("Sub Section added successfully"))
}

// @Summary Add Tag
// @Description Add a new tag to the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param genre_slug path string true "Genre slug"
// @Param data body schemas.TagsAddSchema true "Tag"
// @Success 201 {object} schemas.ResponseSchema "Tag added successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/tags/add/{genre_slug} [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddBookTag(c *fiber.Ctx) error {
	db := ep.DB
	genre := genreManager.GetBySlug(db, c.Params("genre_slug"))
	if genre == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Genre does not exist"))
	}

	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	name := data.Name
	existingTag := models.Tag{}
	db.Where("LOWER(name) = LOWER(?)", name).First(&existingTag)
	if existingTag.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Tag already exists"))
	}
	db.Create(&models.Tag{Name: name, GenreID: genre.ID})
	return c.Status(201).JSON(ResponseMessage("Tag added successfully"))
}

// @Summary Update Genre
// @Description Update a genre.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Genre slug"
// @Param data body schemas.TagsAddSchema true "Genre"
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

	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	name := data.Name
	existingGenre := models.Genre{}
	db.Where("LOWER(name) = LOWER(?)", name).Not("id = ?", genre.ID).First(&existingGenre)
	if existingGenre.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Genre already exists with that name"))
	}
	genre.Name = name
	db.Save(&genre)
	return c.Status(200).JSON(ResponseMessage("Genre updated successfully"))
}

// @Summary Update Section
// @Description Update a section.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Section slug"
// @Param data body schemas.TagsAddSchema true "Section"
// @Success 200 {object} schemas.ResponseSchema "Section Updated Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/sections/{slug} [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateBookSection(c *fiber.Ctx) error {
	db := ep.DB
	section := genreManager.GetSectionBySlug(db, c.Params("slug"))
	if section == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Section does not exist"))
	}

	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	name := data.Name
	existingSection := models.Section{}
	db.Where("LOWER(name) = LOWER(?)", name).Not("id = ?", section.ID).First(&existingSection)
	if existingSection.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "Section already exists with that name"))
	}
	section.Name = name
	db.Save(&section)
	return c.Status(200).JSON(ResponseMessage("Section updated successfully"))
}

// @Summary Update SubSection
// @Description Update a subsection.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "SubSection slug"
// @Param data body schemas.TagsAddSchema true "SubSection"
// @Success 200 {object} schemas.ResponseSchema "SubSection Updated Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/subsections/{slug} [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateBookSubSection(c *fiber.Ctx) error {
	db := ep.DB
	subsection, _ := genreManager.GetSubSectionBySlug(db, c.Params("slug"))
	if subsection == nil {
		return c.Status(404).JSON(utils.NotFoundErr("SubSection does not exist"))
	}

	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	name := data.Name
	existingSubSection := models.SubSection{}
	db.Where("LOWER(name) = LOWER(?)", name).Not("id = ?", subsection.ID).First(&existingSubSection)
	if existingSubSection.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("name", "SubSection already exists with that name"))
	}
	subsection.Name = name
	db.Save(&subsection)
	return c.Status(200).JSON(ResponseMessage("SubSection updated successfully"))
}

// @Summary Update Tag
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
func (ep Endpoint) AdminUpdateBookTag(c *fiber.Ctx) error {
	db := ep.DB

	tag := tagManager.GetBySlug(db, c.Params("slug"))
	if tag == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Tag does not exist"))
	}

	data := schemas.TagsAddSchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
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

// @Summary Delete Genre
// @Description Delete a genre.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Genre slug"
// @Success 200 {object} schemas.ResponseSchema "Genre Deleted Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/genres/{slug} [delete]
// @Security BearerAuth
func (ep Endpoint) AdminDeleteBookGenre(c *fiber.Ctx) error {
	db := ep.DB
	genre := genreManager.GetBySlug(db, c.Params("slug"))
	if genre == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Genre does not exist"))
	}
	db.Model(&genre).Association("Tags").Clear()
	db.Delete(&genre)
	return c.Status(200).JSON(ResponseMessage("Genre deleted successfully"))
}

// @Summary Delete Section
// @Description Delete a book section.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Section slug"
// @Success 200 {object} schemas.ResponseSchema "Section Deleted Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/sections/{slug} [delete]
// @Security BearerAuth
func (ep Endpoint) AdminDeleteBookSection(c *fiber.Ctx) error {
	db := ep.DB
	section := genreManager.GetSectionBySlug(db, c.Params("slug"))
	if section == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Section does not exist"))
	}
	db.Delete(&section)
	return c.Status(200).JSON(ResponseMessage("Section deleted successfully"))
}

// @Summary Delete SubSection
// @Description Delete a book subsection.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "SubSection slug"
// @Success 200 {object} schemas.ResponseSchema "SubSection Deleted Successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/subsections/{slug} [delete]
// @Security BearerAuth
func (ep Endpoint) AdminDeleteBookSubSection(c *fiber.Ctx) error {
	db := ep.DB
	subsection, _ := genreManager.GetSubSectionBySlug(db, c.Params("slug"))
	if subsection == nil {
		return c.Status(404).JSON(utils.NotFoundErr("SubSection does not exist"))
	}
	db.Delete(&subsection)
	return c.Status(200).JSON(ResponseMessage("SubSection deleted successfully"))
}

// @Summary Add Book To A SubSection
// @Description Add a book to a subsection.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "SubSection slug"
// @Param book_slug path string true "Book slug"
// @Success 200 {object} schemas.ResponseSchema "Book added to subsection successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/subsections/{slug}/add-book/{book_slug} [get]
// @Security BearerAuth
func (ep Endpoint) AddBookToSubSection(c *fiber.Ctx) error {
	db := ep.DB

	subsection, _ := genreManager.GetSubSectionBySlug(db, c.Params("slug"))
	if subsection == nil {
		return c.Status(404).JSON(utils.NotFoundErr("SubSection does not exist"))
	}

	book, err := bookManager.GetBySlug(db, c.Params("book_slug"), false)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Check if already exists to prevent duplicates
	var count int64
	db.Model(&models.BookSubSection{}).
		Where("book_id = ? AND sub_section_id = ?", book.ID, subsection.ID).
		Count(&count)

	if count > 0 {
		return c.Status(200).JSON(ResponseMessage("Book already in subsection"))
	}

	// Determine order in section
	var order int64
	db.Model(&models.BookSubSection{}).
		Where("sub_section_id = ?", subsection.ID).
		Count(&order)

	// Add association manually to maintain ordering
	bookSub := models.BookSubSection{
		BookID:         book.ID,
		SubSectionID:   subsection.ID,
		OrderInSection: uint(order + 1),
	}

	if err := db.Create(&bookSub).Error; err != nil {
		return c.Status(500).JSON(utils.ServerErr("Something went wrong"))
	}

	return c.Status(200).JSON(ResponseMessage("Book added to subsection successfully"))
}

// @Summary Remove Book From SubSection
// @Description Remove a book from a subsection and adjust order of remaining books.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "SubSection slug"
// @Param book_slug path string true "Book slug"
// @Success 200 {object} schemas.ResponseSchema "Book removed from subsection successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/subsections/{slug}/remove-book/{book_slug} [get]
// @Security BearerAuth
func (ep Endpoint) RemoveBookFromSubSection(c *fiber.Ctx) error {
	db := ep.DB

	subsection, _ := genreManager.GetSubSectionBySlug(db, c.Params("slug"))
	if subsection == nil {
		return c.Status(404).JSON(utils.NotFoundErr("SubSection does not exist"))
	}

	book, err := bookManager.GetBySlug(db, c.Params("book_slug"), false)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	var bookSub models.BookSubSection
	if err := db.
		Where("book_id = ? AND sub_section_id = ?", book.ID, subsection.ID).
		First(&bookSub).Error; err != nil {
		return c.Status(404).JSON(utils.NotFoundErr("Book not found in subsection"))
	}

	tx := db.Begin()

	// Delete the record
	if err := tx.Delete(&bookSub).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(utils.ServerErr("Failed to remove book from subsection"))
	}

	// Shift order up for other books in that subsection
	if err := tx.Model(&models.BookSubSection{}).
		Where("sub_section_id = ? AND order_in_section > ?", subsection.ID, bookSub.OrderInSection).
		Update("order_in_section", gorm.Expr("order_in_section - 1")).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(utils.ServerErr("Failed to update book order"))
	}

	tx.Commit()

	return c.Status(200).JSON(ResponseMessage("Book removed from subsection successfully"))
}

// @Summary Delete Tag
// @Description Delete a tag from the app.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Tag slug"
// @Success 200 {object} schemas.ResponseSchema "Tag delete successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/tags/{slug} [delete]
// @Security BearerAuth
func (ep Endpoint) AdminDeleteBookTag(c *fiber.Ctx) error {
	db := ep.DB

	tag := tagManager.GetBySlug(db, c.Params("slug"))
	if tag == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Tag does not exist"))
	}
	db.Model(&tag).Association("Genres").Clear()
	db.Delete(&tag)
	return c.Status(200).JSON(ResponseMessage("Tag deleted successfully"))
}

// @Summary Get Sections
// @Description Retrieve sections with sub sections.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Success 200 {object} schemas.SectionsWithSubSectionsSchema "Sections Retrieved Successfully"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/sections [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetSections(c *fiber.Ctx) error {
	db := ep.DB
	sections := genreManager.GetAllSections(db)
	response := schemas.SectionsWithSubSectionsSchema{
		ResponseSchema: ResponseMessage("Sections retrieved successfully"),
	}.Init(sections)
	return c.Status(200).JSON(response)
}

// @Summary Get A Sub Section
// @Description Retrieve a single sub section.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Sub Section slug"
// @Success 200 {object} schemas.SubSectionWithBooksResponseSchema "Sub section Retrieved Successfully"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/subsections/{slug} [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetSubSection(c *fiber.Ctx) error {
	db := ep.DB
	subSection, books := genreManager.GetSubSectionBySlug(db, c.Params("slug"))
	if subSection == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Sub section does not exist"))
	}
	// Paginate books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 200)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	castedBooks := paginatedBooks.([]schemas.BookWithOrder)

	response := schemas.SubSectionWithBooksResponseSchema{
		ResponseSchema: ResponseMessage("Sub section retrieved successfully"),
		Data: schemas.SubSectionWithBooksSchema{}.Init(*subSection, castedBooks, *paginatedData),
	}
	return c.Status(200).JSON(response)
}

// @Summary List Books with Pagination
// @Description Retrieves a list of books with support for pagination and optional filtering based on book title.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param page query int false "Current Page" default(1)
// @Param title query string false "Title of the book to filter by"
// @Param name query string false "name or username of the book author to filter by"
// @Param rating query bool false "Filter by highest ratings"
// @Param genre_slug query string false "Filter by Genre slug"
// @Param section_slug query string false "Filter by Section slug"
// @Param sub_section_slug query string false "Filter by Sub Section slug"
// @Param tag_slug query string false "Filter by Tag slug"
// @Param featured query bool false "Filter by Featured"
// @Param weeklyFeatured query bool false "Filter by Weekly Featured"
// @Param trending query bool false "Filter by Trending"
// @Success 200 {object} schemas.BooksResponseSchema "Successfully retrieved list of books"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetBooks(c *fiber.Ctx) error {
	db := ep.DB
	titleQuery := c.Query("title", "")
	ratingQuery := c.QueryBool("rating", false)
	nameQuery := c.Query("name", "")
	genreSlug := c.Query("genre_slug", "")
	sectionSlug := c.Query("section_slug", "")
	subSectionSlug := c.Query("sub_section_slug", "")
	tagSlug := c.Query("tag_slug", "")
	featured := c.QueryBool("featured")
	weeklyFeatured := c.QueryBool("weekly_featured")
	trending := c.QueryBool("trending")

	books, _ := bookManager.GetLatest(db, genreSlug, sectionSlug, subSectionSlug, tagSlug, titleQuery, ratingQuery, "", nameQuery, featured, weeklyFeatured, trending, false)

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

// @Summary List Author Books with Pagination
// @Description Retrieves a list of a particular author books with support for pagination and optional filtering based on book title.
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param username path string true "Username of the author"
// @Param page query int false "Current Page" default(1)
// @Param title query string false "Title of the book to filter by"
// @Param rating query bool false "Filter by highest ratings"
// @Param genre_slug query string false "Filter by Genre slug"
// @Param section_slug query string false "Filter by Section slug"
// @Param sub_section_slug query string false "Filter by Sub Section slug"
// @Param tag_slug query string false "Filter by Tag slug"
// @Param featured query bool false "Filter by Featured"
// @Param weeklyFeatured query bool false "Filter by Weekly Featured"
// @Param trending query bool false "Filter by Trending"
// @Success 200 {object} schemas.BooksResponseSchema "Successfully retrieved list of books"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/by-username/{username} [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetAuthorBooks(c *fiber.Ctx) error {
	db := ep.DB
	titleQuery := c.Query("title", "")
	ratingQuery := c.QueryBool("rating", false)
	username := c.Params("username")
	genreSlug := c.Query("genre_slug", "")
	sectionSlug := c.Query("section_slug", "")
	subSectionSlug := c.Query("sub_section_slug", "")
	tagSlug := c.Query("tag_slug", "")
	featured := c.QueryBool("featured")
	weeklyFeatured := c.QueryBool("weekly_featured")
	trending := c.QueryBool("trending")

	author := models.User{Username: username, AccountType: choices.ACCTYPE_AUTHOR}
	db.Take(&author, author)

	if author.ID == uuid.Nil {
		return c.Status(404).JSON(utils.NotFoundErr("Author does not exist!"))
	}

	books, _ := bookManager.GetLatest(db, genreSlug, sectionSlug, subSectionSlug, tagSlug, titleQuery, ratingQuery, username, "", featured, weeklyFeatured, trending, false)

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

// @Summary View Book Details
// @Description This endpoint allows an admin to view details of a book
// @Tags Admin | Books
// @Param page query int false "Current Page (for reviews pagination)" default(1)
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.BookDetailResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /admin/books/book-detail/{slug} [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetBookDetails(c *fiber.Ctx) error {
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

// @Summary View Book Reading Progress Data
// @Description This endpoint allows an admin to view reading progress data of a book
// @Tags Admin | Books
// @Param slug path string true "Book slug"
// @Param time_range_filter query string false "Time Range to filter by (12_months, 30_days, 7_days)" default(12_months)
// @Success 200 {object} schemas.BookReadingProgressResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /admin/books/book-detail/{slug}/reading-progress [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetBookReadingProgress(c *fiber.Ctx) error {
	db := ep.DB
	time_range_filter := c.Query("time_range_filter", "12_months")
	if time_range_filter != "12_months" && time_range_filter != "30_days" && time_range_filter != "7_days" {
		return c.Status(400).JSON(utils.InvalidParamErr("Invalid date filter"))
	}
	book, err := bookManager.GetBySlug(db, c.Params("slug"), false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	var readingProgressData []schemas.BookReadingProgressSchema
	switch time_range_filter {
	case "30_days":
		readingProgressData = bookManager.Get30DayReadingProgress(ep.DB, book.ID)
	case "7_days":
		readingProgressData = bookManager.Get7DayReadingProgress(ep.DB, book.ID)
	default:
		readingProgressData = bookManager.GetYearlyReadingProgress(ep.DB, book.ID)
	}
	response := schemas.BookReadingProgressResponseSchema{
		ResponseSchema: ResponseMessage("Book Reading Progress fetched successfully"),
		Data:           readingProgressData,
	}
	return c.Status(200).JSON(response)
}

// @Summary View Book Retention Data
// @Description This endpoint allows an admin to view retention data of a book
// @Tags Admin | Books
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.BookRetentionStatsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /admin/books/book-detail/{slug}/retention-stats [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetBookRetentionStats(c *fiber.Ctx) error {
	db := ep.DB
	book, err := bookManager.GetBySlug(db, c.Params("slug"), false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	retentionData := bookManager.GetReaderRetentionPieData(db, book.ID)
	response := schemas.BookRetentionStatsResponseSchema{
		ResponseSchema: ResponseMessage("Book Retention data fetched successfully"),
		Data:           retentionData,
	}
	return c.Status(200).JSON(response)
}

// @Summary Toggle Book Completed Status
// @Description Set the book completed status to true or false. 
// @Tags Admin | Books
// @Accept json
// @Produce json
// @Param slug path string true "Book slug"
// @Success 200 {object} schemas.BookCompletionStatusResponseSchema "Book status updated successfully"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books/book/{slug}/toggle-book-completion-status [get]
// @Security BearerAuth
func (ep Endpoint) ToggleBookCompletionStatus(c *fiber.Ctx) error {
	db := ep.DB
	book, err := bookManager.GetBySlug(db, c.Params("slug"), false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	if book.Completed {
		book.Completed = false
	} else {
		book.Completed = true
	}
	db.Save(&book)
	response := schemas.BookCompletionStatusResponseSchema{
		ResponseSchema: ResponseMessage("Book completed to subsection successfully"),
		Data: schemas.BookCompletionStatusSchema{Completed: book.Completed},
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
		return c.Status(400).JSON(utils.InvalidParamErr("Invalid contract status"))
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
