package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
)

var truthy = true
// @Summary List Users with Pagination
// @Description Retrieves a list of user profiles with support for pagination and optional filtering based on user account type.
// @Tags Admin
// @Accept json
// @Produce json
// @Param account_type query string false "Type of user to filter by" Enums(READER, WRITER, ADMIN)
// @Param page query int false "Current page" default(1)
// @Success 200 {object} schemas.UserProfilesResponseSchema "Successfully retrieved list of user profiles"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/users [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetUsers(c *fiber.Ctx) error {
	db := ep.DB
	acctType := c.Query("account_type", "")
	var accountType *choices.AccType
	var staff *bool

	if acctType == "" {
		accountType = nil
	} else {
		if !IsAmongUserType(acctType) {
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid account type"))
		}
		if acctType == "ADMIN" {
			staff = &truthy
		} else {
			acc := choices.AccType(acctType)
			accountType = &acc
		}
	} 
	users := userManager.GetAll(db, accountType, staff)
	// Paginate and return users
	paginatedData, paginatedUsers, err := PaginateQueryset(users, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	users = paginatedUsers.([]models.User)

	response := schemas.UserProfilesResponseSchema{
		ResponseSchema: ResponseMessage("Profiles fetched successfully"),
		Data: schemas.UserProfilesResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(users),
	}

	return c.Status(200).JSON(response)
}

// @Summary Update User Role
// @Description Updates the account type of a specified user.
// @Tags Admin
// @Accept json
// @Produce json
// @Param data body schemas.UpdateUserRoleSchema true "User role update data"
// @Success 200 {object} schemas.UserProfileResponseSchema "Successfully updated user details"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/users/user [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateUser(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.UpdateUserRoleSchema{}

	if errCode, errData := ValidateRequest(c, &data);

	errData != nil{
		return c.Status(*errCode).JSON(errData)
	}

	account_type := choices.AccType(data.AccountType)

	var user models.User

	if data.Username != nil{
		
		result := db.Where("username = ?", data.Username).First(&user)
		if result.Error != nil{
			return c.Status(404).JSON(utils.RequestErr(utils.ERR_NOT_FOUND, "User Not Found"))
		}

		user.AccountType = account_type
	}

	db.Save(&user)

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("User details updated successfully!"),
		Data: schemas.UserProfile{}.Init(user),
	}
	return c.Status(200).JSON(response)
}