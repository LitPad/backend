package helpers

import (
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserNotFoundError(c *fiber.Ctx, msg string, err error) error {
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(utils.NotFoundErr(msg))
	}
	return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Internal server error"))
}
