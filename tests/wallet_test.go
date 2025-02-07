package tests

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getAvailbleCoins(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	TestCoin(db)

	t.Run("Accept Coins Fetch", func(t *testing.T) {
		url := baseUrl + "/coins"
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Coins fetched successfully", body["message"])
	})
}

// func buyCoins(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
// 	coin := TestCoin(db)
// 	coinData := schemas.BuyCoinSchema{CoinID: uuid.New(), Quantity: 1}

// 	t.Run("Reject Coins Buy Due To Invalid Coin ID", func(t *testing.T) {
// 		url := baseUrl + "/coins"
// 		res := ProcessJsonTestBody(t, app, url, "POST", coinData, token)
// 		// Assert Status code
// 		assert.Equal(t, 404, res.StatusCode)

// 		// Parse and assert body
// 		body := ParseResponseBody(t, res.Body).(map[string]interface{})
// 		assert.Equal(t, "failure", body["status"])
// 		assert.Equal(t, "No set of coins with that ID", body["message"])
// 	})

// 	t.Run("Accept Coins Buying Due To Valid Data", func(t *testing.T) {
// 		coinData.CoinID = coin.ID
// 		url := baseUrl + "/coins"
// 		res := ProcessJsonTestBody(t, app, url, "POST", coinData, token)
// 		// Assert Status code
// 		assert.Equal(t, 201, res.StatusCode)

// 		// Parse and assert body
// 		body := ParseResponseBody(t, res.Body).(map[string]interface{})
// 		assert.Equal(t, "success", body["status"])
// 		assert.Equal(t, "Payment Data Generated", body["message"])
// 	})
// }

func TestWallet(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	user := TestVerifiedUser(db)
	token := AccessToken(db, user)
	baseUrl := "/api/v1/wallet"

	// Run Wallet Endpoint Tests
	getAvailbleCoins(t, app, db, baseUrl, token)
	// buyCoins(t, app, db, baseUrl, token)
}