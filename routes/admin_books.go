package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
)

// @Summary List Books with Pagination
// @Description Retrieves a list of books with support for pagination and optional filtering based on book title.
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Current Page" default(1)
// @Param title query string false "Title of the book to filter by"
// @Success 200 {object} schemas.BooksResponseSchema "Successfully retrieved list of books"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/books [get]
func (ep Endpoint) AdminGetBooks(c *fiber.Ctx) error {
	db := ep.DB

	titleQuery := c.Query("title", "")

	books := []models.Book{}
	query := db
	if len(titleQuery) > 0 {
		query = query.Where("title ILIKE ?", "%"+titleQuery+"%")
	}
	query.Find(&books)

	// Paginate and return books
	paginatedData, paginatedBooks, err := PaginateQueryset(books, c, 400)
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
