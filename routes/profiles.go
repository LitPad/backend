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
// @Tags User
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/{username} [get]
func (ep Endpoint) GetProfile(c *fiber.Ctx) error {
	db := ep.DB

	pathParams := c.Params("username")

	if pathParams == ""{
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path params"))
	}

	user := models.User{Username: pathParams}
	db.Take(&user, user)

	if user.ID == uuid.Nil || !user.IsEmailVerified {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "Cannot Find Resource"))
	}
	
	response := schemas.UserProfileResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Profile fetched successfully"}.Init(),
		Data: schemas.UserProfile{
			FirstName: user.FirstName,
			LastName: user.LastName,
			Username: user.Username,
			Email: user.Email,
			Avatar: user.Avatar,
			Bio: user.Bio,
			AccountType: schemas.AccType(user.AccountType),
		},
	}

	return c.Status(200).JSON(response)
}

// @Summary Update User Profile
// @Description This endpoint updates a user's profile
// @Tags User
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update [patch]
func (ep Endpoint) UpdateProfile(c *fiber.Ctx)error {
	db := ep.DB

	data := schemas.UpdateUserProfileSchema{}

	if errCode, errData := ValidateRequest(c, &data); errData !=nil{
		return c.Status(*errCode).JSON(errData)
	}

	savedUser := c.Locals("user").(*models.User)

	if savedUser == nil || savedUser.ID == uuid.Nil{
		return c.Status(403).JSON(utils.RequestErr(utils.ERR_UNAUTHORIZED_USER, "SignIn to make this request"))
	}

	if len(data.Username) > 0{
	searchUser := models.User{Username: data.Username}

		db.Take(&searchUser, searchUser)
		if searchUser.ID != uuid.Nil{
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Username is already taken by another user"))
		}
	}

	// current design supports update of username only hence the code looking like this

	db.Model(&models.User{}).Where("username = ?", savedUser.Username).Update("username", data.Username)

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "User details updated successfully"}.Init(),
		Data: schemas.UserProfile{
			Username: data.Username,
			Email: savedUser.Email,
		},
	}

	return c.Status(200).JSON(response)
}

// @Summary Update User Password
// @Description This endpoint updates a user's password
// @Tags User
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update-password [put]
func (ep Endpoint) UpdatePassword(c *fiber.Ctx)error {
	db := ep.DB

	data := schemas.UpdatePasswordSchema{}

	if errCode, errData := ValidateRequest(c, &data); errData !=nil{
		return c.Status(*errCode).JSON(errData)
	}

	user := c.Locals("user").(*models.User);

	searchUserInterface := models.User{Email: user.Email}

	db.Take(&searchUserInterface, searchUserInterface)

	if searchUserInterface.ID == uuid.Nil{
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "cannot find user"))
	}

	if !utils.CheckPasswordHash(data.OldPassword, searchUserInterface.Password) {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_PASSWORD_MISMATCH, "Password Mismatch"))
	}

	if utils.CheckPasswordHash(data.NewPassword, searchUserInterface.Password){
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_PASSWORD_SAME, "new password is same as old password"))
	}

	searchUserInterface.Password = utils.HashPassword(data.NewPassword)
	db.Save(&searchUserInterface)

	response := schemas.ResponseSchema{Message: "Password updated successfully"}.Init()
	
	return c.Status(200).JSON(response)
}
