package routes

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
	if !user.IsStaff {
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

	if len(token) < 1 {
		return c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Forbidden"))
	}

	if !strings.HasPrefix(token, "Litpad ") {
		return c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Forbidden"))
	}

	parsedToken, err := jwt.Parse(token[7:], func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.WalletSecret), nil
	})

	if err != nil {
		c.Status(500).JSON(utils.ERR_SERVER_ERROR)
	}

	if !parsedToken.Valid {
		c.Status(403).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "Forbidden"))
	}

	return c.Next()
}

func (ep Endpoint) DynamicRateLimiter(expirationMinute int, maxRequest int ) fiber.Handler {
	// Apply rate limiter middleware
	return limiter.New(limiter.Config{
		Max:        maxRequest,               // dynamic request limit
		Expiration: time.Duration(expirationMinute) * time.Minute, // expiration time in minutes
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Use IP as the key for rate limiting
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(utils.RateLimitError("Rate Limit Reached"))
		},
	})
}

func RequestLogger(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !strings.HasPrefix(c.Path(), "/api/v1") {
			return c.Next()
		}

		// Call the next middleware to capture the response status
		err := c.Next()

		// Initialize bodyStr as empty
		bodyStr := ""

		// Get the Content-Type
		contentType := c.Get("Content-Type")

		// Handle multipart/form-data separately
		if strings.HasPrefix(contentType, "multipart/form-data") {
			// For multipart form data, log form fields and file metadata
			form, err := c.MultipartForm()
			if err != nil {
				log.Println("Error reading multipart form:", err)
			} else {
				// Extract form fields (key-value pairs)
				formData := make(map[string]string)
				for key, values := range form.Value {
					formData[key] = strings.Join(values, ",") // Join values if there are multiple for a key
				}

				// Log file metadata (size, filename, etc.)
				for key, files := range form.File {
					for _, file := range files {
						log.Printf("Received file: %s, Size: %d bytes, Field: %s", file.Filename, file.Size, key)
					}
				}

				// Marshal form data into a string for logging
				parsedBody, _ := json.Marshal(formData)
				bodyStr = string(parsedBody)
			}
		} else if strings.HasPrefix(contentType, "application/json") {
			// For JSON, we can safely unmarshal
			bodyBytes := c.Body()
			if len(bodyBytes) > 0 {
				var bodyMap map[string]interface{}
				if json.Unmarshal(bodyBytes, &bodyMap) == nil { // Valid JSON
					parsedBody, _ := json.Marshal(bodyMap)
					bodyStr = string(parsedBody)
				} else {
					bodyStr = string(bodyBytes) // Fallback for invalid JSON
				}
			}
		} else {
			// Handle other content types as a fallback (plain text, etc.)
			bodyBytes := c.Body()
			if len(bodyBytes) > 0 {
				bodyStr = string(bodyBytes)
			}
		}

		// Create the log entry
		log := models.Log{
			Method:     c.Method(),
			Path:       c.Path(),
			IP:         c.IP(),
			StatusCode: c.Response().StatusCode(),
			QueryParams: string(c.Request().URI().QueryString()),
			PathParams: string(c.Params("*")),
			Body:       bodyStr, // Log the parsed body (or form data)
		}

		// Save the log
		db.Create(&log)
		return err
	}
}


func ParseUUID(input string) *uuid.UUID {
	uuidVal, err := uuid.Parse(input)
	if err != nil {
		return nil
	}
	return &uuidVal
}
