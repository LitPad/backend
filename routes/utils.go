package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
)

func ResponseMessage(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}

func RequestUser(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}
