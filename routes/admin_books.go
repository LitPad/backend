package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
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
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Current Page" default(1)
// @Param name query string false "Name of the author to filter by"
// @Param contract_status query string false "status of the contract to filter by" Enums(PENDING, APPROVED, DECLINED, UPDATED)
// @Success 200 {object} schemas.ContractsResponseSchema "Successfully retrieved list of book contracts"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/contracts [get]
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
