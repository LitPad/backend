package routes

import (
	"strings"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUser(token string, db *gorm.DB) (*models.User, *string) {
	if !strings.HasPrefix(token, "Bearer ") {
		err := "Auth Bearer Not Provided!"
		return nil, &err
	}
	user, err := DecodeAccessToken(token[7:], db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ep Endpoint) AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB

	if len(token) < 1 {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNAUTHORIZED_USER, "Unauthorized User!"))
	}
	user, err := GetUser(token, db)
	if err != nil {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, *err))
	}
	c.Locals("user", user)
	return c.Next()
}

func (ep Endpoint) AuthOrGuestMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB
	c.Locals("user", &models.User{})
	if len(token) > 1 {
		user, err := GetUser(token, db)
		if err != nil {
			return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, *err))
		}
		c.Locals("user", user)
	}
	return c.Next()
}

func (ep Endpoint) AuthorMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB

	if len(token) < 1 {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNAUTHORIZED_USER, "Unauthorized User!"))
	}
	user, err := GetUser(token, db)
	if err != nil {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, *err))
	}
	if user.AccountType != choices.ACCTYPE_WRITER {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_AUTHORS_ONLY, "For Authors only!"))
	}
	c.Locals("user", user)
	return c.Next()
}

func (ep Endpoint) AdminMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB

	if len(token) < 1 {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNAUTHORIZED_USER, "Unauthorized User!"))
	}
	user, err := GetUser(token, db)
	if err != nil {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, *err))
	}
	if !user.IsStaff {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_ADMINS_ONLY, "For Admin only!"))
	}
	c.Locals("user", user)
	return c.Next()
}

func (ep Endpoint) WalletAccessMiddleware(c *fiber.Ctx) error {

	conf := config.GetConfig()

	token := c.Get("Access")

	if len(token) < 1{
		return c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Forbidden"))
	}

	if !strings.HasPrefix(token, "Litpad "){
		return c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Forbidden"))
	}

	parsedToken, err := jwt.Parse(token[7:], func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.WalletSecret), nil
	})

	if err != nil{
		c.Status(500).JSON(utils.ERR_SERVER_ERROR)
	}

	if !parsedToken.Valid{
		c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Forbidden"))
	}

	return c.Next()
}

func ParseUUID(input string) *uuid.UUID {
    uuidVal, err := uuid.Parse(input)
	if err != nil {
		return nil
	}
    return &uuidVal
}