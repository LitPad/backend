package routes

import (
	"github.com/gofiber/fiber/v2"
)

// @Summary View User Profile
// @Description This endpoint views a user profile
// @Tags Auth
// @Success 200 {object} schemas.ResponseSchema
// @Failure 401 {object} utils.ErrorResponse
// @Router /profiles/{username} [get]
func (ep Endpoint) GetProfile(c *fiber.Ctx) error {
	// db := ep.DB
	return c.Status(200).JSON("")
}