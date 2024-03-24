package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary View User Profile
// @Description This endpoint views a user profile
// @Tags Profiles
// @Param username path string true "Username of user"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/profile/{username} [get]
func (ep Endpoint) GetProfile(c *fiber.Ctx) error {
	db := ep.DB

	pathParams := c.Params("username")

	if pathParams == "" {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path params"))
	}

	user := models.User{Username: pathParams, IsEmailVerified: true}
	db.Take(&user, user)

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "Cannot Find Resource"))
	}

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("Profile fetched successfully"),
		Data:           schemas.UserProfile{}.Init(user),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update User Profile
// @Description This endpoint updates a user's profile
// @Tags Profiles
// @Param profile body schemas.UpdateUserProfileSchema true "Profile object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update [patch]
// @Security BearerAuth
func (ep Endpoint) UpdateProfile(c *fiber.Ctx) error {
	db := ep.DB
	savedUser := RequestUser(c)
	data := schemas.UpdateUserProfileSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	username := *data.Username
	if data.Username != nil {
		searchUser := models.User{Username: username}
		db.Not(models.User{BaseModel: models.BaseModel{ID: savedUser.ID}}).Take(&searchUser, searchUser)
		if searchUser.ID != uuid.Nil {
			data := map[string]string{"username": "Username is already taken by another user"}
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid Entry", data))
		}
		savedUser.Username = username
	}

	// current design supports update of username only hence the code looking like this
	db.Save(&savedUser)

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("User details updated successfully"),
		Data:           schemas.UserProfile{}.Init(*savedUser),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update User Password
// @Description This endpoint updates a user's password
// @Tags Profiles
// @Param profile body schemas.UpdatePasswordSchema true "Password object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update-password [put]
// @Security BearerAuth
func (ep Endpoint) UpdatePassword(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.UpdatePasswordSchema{}

	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := RequestUser(c)

	if !utils.CheckPasswordHash(data.OldPassword, user.Password) {
		data := map[string]string{"old_password": "Password Mismatch"}

		return c.Status(400).JSON(utils.RequestErr(utils.ERR_PASSWORD_MISMATCH, "Invalid Entry", data))
	}

	if data.NewPassword == data.OldPassword {
		data := map[string]string{"new_password": "new password is same as old password"}
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_PASSWORD_SAME, "Invalid Entry", data))
	}

	user.Password = utils.HashPassword(data.NewPassword)
	// Clear tokens to logout user
	user.Access = nil
	user.Refresh = nil
	db.Save(&user)
	return c.Status(200).JSON(ResponseMessage("Password updated successfully"))
}
