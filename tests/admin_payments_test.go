package tests

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getTransactions(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	t.Run("Accept Transactions Fetch", func(t *testing.T) {
		TestSubscriptionPlan(db)
		url := fmt.Sprintf("%s/transactions", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Transactions fetched successfully", body["message"])
	})
}

func TestAdminPayments(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	admin := TestAdmin(db)
	token := AccessToken(db, admin)
	baseUrl := "/api/v1/admin/payments"

	// Run Admin Payments Endpoint Tests
	getTransactions(t, app, db, baseUrl, token)
}
