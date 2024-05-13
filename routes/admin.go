package routes

import (
	"strconv"
	"strings"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// endpoint that returns all books

// endpoint that returns all transactions

// endpoint that returns all wallets

// @Summary List Users with Pagination
// @Description Retrieves a list of user profiles with support for pagination and optional filtering based on user account type.
// @Tags Users
// @Accept json
// @Produce json
// @Param type query string false "Type of user to filter by (all, reader, writer)" Enums(all, reader, writer)
// @Param limit query int false "Limit number of user profiles per page (default is 10)" default(10)
// @Param page query int false "Page number starting from 0 (default is 0)" default(0)
// @Success 200 {object} schemas.UserProfilesResponseSchema "Successfully retrieved list of user profiles"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 404 {object} utils.ErrorResponse "No users found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/users [get]
func (ep Endpoint) GetUsers(c *fiber.Ctx) error{
	db := ep.DB

	limitQuery := c.Query("limit", "10")
	pageQuery := c.Query("page", "0")
	userType := c.Query("type", "all")

	limit, err := strconv.Atoi(limitQuery)

	if err !=nil{
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid query param `limit` "))
	}

	page, err := strconv.Atoi(pageQuery)

	if err !=nil{
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid query param `page` "))
	}

	offset := (page - 1) * limit

	var users []models.User

	query := db.Scopes(scopes.FollowerFollowingPreloaderScope)
	
	if userType == strings.ToLower(string(choices.ACCTYPE_READER)) {
        query = query.Where("account_type = ?", "READER")
		
    } else if userType == strings.ToLower(string(choices.ACCTYPE_WRITER)) {
		
        query = query.Where("account_type = ?", "WRITER")
    }

	if err = query.Offset(offset).Limit(limit).Find(&users).Error; err != nil{
		if err == gorm.ErrRecordNotFound{
			response := schemas.UserProfilesResponseSchema{
				ResponseSchema: ResponseMessage("No profiles exist"),
				Data: []schemas.UserProfile{},
			}

			return c.Status(200).JSON(response)
		}

		return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Internal Server Error"))
	}

	profiles := make([]schemas.UserProfile, len(users))

	for i, user := range users{
		profiles[i] = schemas.UserProfile{}.Init(user)
	}

	response := schemas.UserProfilesResponseSchema{
		ResponseSchema: ResponseMessage("Profiles fetched successfully"),
		Data: profiles,
	}

	return c.Status(200).JSON(response)
}


// @Summary List Books with Pagination
// @Description Retrieves a list of books with support for pagination and optional filtering based on book title.
// @Tags Books
// @Accept json
// @Produce json
// @Param title query string false "Title of the book to filter by"
// @Param limit query int false "Limit number of book profiles per page (default is 10)" default(10)
// @Param page query int false "Page number starting from 0 (default is 0)" default(0)
// @Success 200 {object} schemas.BookResponseSchema "Successfully retrieved list of books"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 404 {object} utils.ErrorResponse "No books found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /books [get]
func (ep Endpoint) GetBooks(c *fiber.Ctx) error{
	db := ep.DB

	limitQuery := c.Query("limit", "10")
	pageQuery := c.Query("page", "0")
	titleQuery := c.Query("title", "")

	limit, err := strconv.Atoi(limitQuery)

	if err !=nil{
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid query param `limit` "))
	}

	page, err := strconv.Atoi(pageQuery)

	if err !=nil{
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid query param `page` "))
	}

	offset := (page - 1) * limit

	var books []models.Book

	query := db.Model(&models.Book{})

	if len(titleQuery) > 0 {
		query = query.Where("title ILIKE ?", "%" + titleQuery + "%")
	}
		
	if err = query.Offset(offset).Limit(limit).Find(&books).Error; err != nil{
		if err == gorm.ErrRecordNotFound{
			response := schemas.BookResponseSchema{
				ResponseSchema: ResponseMessage("No books exist"),
				Data: []schemas.BookSchema{},
			}

			return c.Status(200).JSON(response)
		}

		return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Internal Server Error"))
	}

	booksArr := make([]schemas.BookSchema, len(books))

	for i, book := range books{
		booksArr[i] = schemas.BookSchema{}.Init(book)
	}

	response := schemas.BookResponseSchema{
		ResponseSchema: ResponseMessage("Books fetched successfully"),
		Data: booksArr,
	}

	return c.Status(200).JSON(response)
}