package tests

import (
	"fmt"
	"testing"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getUsers(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	t.Run("Reject Users Fetch Due to Invalid Account Type", func(t *testing.T) {
		url := fmt.Sprintf("%s?account_type=invalid-account-type", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid account type!", body["message"])
	})

	t.Run("Accept Users Fetch", func(t *testing.T) {
		TestVerifiedUser(db)
		res := ProcessTestGetOrDelete(app, baseUrl, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Profiles fetched successfully", body["message"])
	})
}

func updateUserRole(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, admin models.User, token string) {
	data := schemas.UpdateUserRoleSchema{AccountType: choices.ACCTYPE_READER}

	t.Run("Reject Role Update Due to Invalid Username", func(t *testing.T) {
		url := fmt.Sprintf("%s/invalid-username", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "PUT", data, token)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "User Not Found", body["message"])
	})

	t.Run("Accept Role Update Due to Valid Username", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", baseUrl, admin.Username)
		res := ProcessJsonTestBody(t, app, url, "PUT", data, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User details updated successfully!", body["message"])
	})
}

func toggleUserActivation(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	user := TestUser(db)

	t.Run("Accept User Activation Toggle", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s/toggle-activation", baseUrl, user.Username)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User deactivated successfully", body["message"])
	})
}

func TestAdminUsers(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	admin := TestAdmin(db)
	token := AccessToken(db, admin)
	baseUrl := "/api/v1/admin/users"

	// Run Admin Users Endpoint Tests
	getUsers(t, app, db, baseUrl, token)
	updateUserRole(t, app, db, baseUrl, admin, token)
	toggleUserActivation(t, app, db, baseUrl, token)
}
