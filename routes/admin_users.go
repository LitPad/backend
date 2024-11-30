package routes

import (
	"fmt"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var truthy = true
// @Summary List Users with Pagination
// @Description Retrieves a list of user profiles with support for pagination and optional filtering based on user account type.
// @Tags Admin | Users
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
// @Tags Admin | Users
// @Accept json
// @Produce json
// @Param username path string true "Username" default(username)
// @Param data body schemas.UpdateUserRoleSchema true "User role update data"
// @Success 200 {object} schemas.UserProfileResponseSchema "Successfully updated user details"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/users/{username} [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateUser(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.UpdateUserRoleSchema{}
	errCode, errData := ValidateRequest(c, &data);
	if errData != nil{
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Username: c.Params("username")}
	db.Scopes(scopes.FollowerFollowingPreloaderScope).Take(&user, user)
	if user.ID == uuid.Nil{
		return c.Status(404).JSON(utils.NotFoundErr("User Not Found"))
	}
	user.AccountType = data.AccountType
	db.Save(&user)

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("User details updated successfully!"),
		Data: schemas.UserProfile{}.Init(user),
	}
	return c.Status(200).JSON(response)
}

// @Summary Reactivate/Deactivate User
// @Description Allows the admin to deactivate/reactivate a user.
// @Tags Admin | Users
// @Param username path string true "Username" default(username)
// @Accept json
// @Produce json
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/users/{username}/toggle-activation [get]
// @Security BearerAuth
func (ep Endpoint) ToggleUserActivation(c *fiber.Ctx) error {
	db := ep.DB
	username := c.Params("username")
	user := models.User{Username: username}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.NotFoundErr("User with that username not found"))
	}
	responseMessageSubstring := "deactivated"
	if user.IsActive {
		user.IsActive = false
	} else {
		responseMessageSubstring = "reactivated"
		user.IsActive = true
	}
	db.Save(&user)
	return c.Status(200).JSON(ResponseMessage(fmt.Sprintf("User %s successfully", responseMessageSubstring)))
}

func (ep Endpoint) AdminGetWaitlist(c *fiber.Ctx)error{
	db := ep.DB
	
	var waitlist []models.Waitlist

	// Preload the Genre details for each waitlist entry

	if err := db.Preload("Genre").Find(&waitlist).Error; err != nil {
		return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to retrieve waitlist"))
	}

	paginatedData, paginatedWaitlist, err := PaginateQueryset(waitlist,c, 100)

	if err != nil{
		return c.Status(400).JSON(err)
	}

	waitlist = paginatedWaitlist.([]models.Waitlist)

	response := schemas.WaitlistListResponseSchema{
		ResponseSchema: ResponseMessage("Waitlist fetched successfully"),
		Data: schemas.WaitlistResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(waitlist),
	}

	return c.Status(200).JSON(response)
}