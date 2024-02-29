package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/senders"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Register a new user
// @Description This endpoint registers new users into our application.
// @Tags Auth
// @Param user body schemas.RegisterUser true "User data"
// @Success 201 {object} schemas.RegisterResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (ep Endpoint) Register(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.RegisterUser{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := utils.ConvertStructData(data, models.User{}).(*models.User)
	db.Take(&user, models.User{Email: user.Email})
	if user.ID != uuid.Nil {
		data := map[string]string{
			"email": "Email already registered!",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Create User
	db.Create(&user)

	// Send Email
	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	db.Save(&otp) // Create or save
	go senders.SendEmail(user, "activate", otp.Code)

	response := schemas.RegisterResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Registration successful"}.Init(),
		Data:           schemas.EmailRequestSchema{Email: user.Email},
	}
	return c.Status(201).JSON(response)
}
