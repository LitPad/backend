package tests

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getSubscribers(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	t.Run("Reject Subscribers Fetch Due to Invalid Subscription Type", func(t *testing.T) {
		url := fmt.Sprintf("%s/?sub_type=invalid-sub-type", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid subscription type!", body["message"])
	})

	t.Run("Reject Subscribers Fetch Due to Invalid Subscription Status", func(t *testing.T) {
		url := fmt.Sprintf("%s/?sub_status=invalid-sub-status", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid subscription status!", body["message"])
	})

	t.Run("Accept Subscribers Fetch Due to Valid Data", func(t *testing.T) {
		TestSubscriber(db)
		res := ProcessTestGetOrDelete(app, baseUrl, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Subscribers fetched successfully", body["message"])
	})
}

func TestAdminSubscribers(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	admin := TestAdmin(db)
	token := AccessToken(db, admin)
	baseUrl := "/api/v1/admin/subscribers"

	// Run Admin Subscribers Endpoint Tests
	getSubscribers(t, app, db, baseUrl, token)
}
