package tests

import (
	"fmt"
	"testing"

	"github.com/LitPad/backend/database"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func register(t *testing.T, app *fiber.App, baseUrl string) {
	t.Run("Register User Successfully", func(t *testing.T) {
		url := fmt.Sprintf("%s/register", baseUrl)
		validEmail := "testregisteruser@email.com"
		userData := schemas.RegisterUser{
			Email:    validEmail,
			Password: "testregisteruserpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", userData)

		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Registration successful", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})

	t.Run("Register User Failure By Already Used details", func(t *testing.T) {
		email := "testregisteruser@email.com"
		url := fmt.Sprintf("%s/register", baseUrl)
		userData := schemas.RegisterUser{
			Email:    email,
			Password: "testregisteruserpassword",
		}

		// Verify that a user with the same email cannot be registered again
		res := ProcessTestBody(t, app, url, "POST", userData)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_ENTRY, body["code"])
		assert.Equal(t, "Invalid Entry", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = "Email already taken!"
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})
}

func TestAuth(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	baseUrl := "/api/v1/auth"

	// Run Auth Endpoint Tests
	register(t, app, baseUrl)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
