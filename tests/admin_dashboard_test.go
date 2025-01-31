package tests

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func getDashboardData(t *testing.T, app *fiber.App, baseUrl string, token string) {
	t.Run("Reject Dashboard Data Fetch Due to Invalid Growth Filter", func(t *testing.T) {
		url := fmt.Sprintf("%s/?user_growth_filter=1", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid user growth filter choice!", body["message"])
	})

	t.Run("Accept Dashboard Data Fetch Due to Valid Data", func(t *testing.T) {
		res := ProcessTestGetOrDelete(app, baseUrl, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Admin Dashboard data retrieved successfully", body["message"])
	})
}

func TestAdminDashboard(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	admin := TestAdmin(db)
	token := AccessToken(db, admin)
	baseUrl := "/api/v1/admin"

	// Run Admin Dashboard Endpoint Tests
	getDashboardData(t, app, baseUrl, token)
}
