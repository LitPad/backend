package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
)

// @Summary View Available Coins
// @Description This endpoint returns all available coins for sale
// @Tags Wallet
// @Success 200 {object} schemas.CoinsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /wallet/coins [get]
func (ep Endpoint) AvailableCoins(c *fiber.Ctx) error {
	db := ep.DB
	coins := []models.Coin{}
	db.Find(&coins)
	response := schemas.CoinsResponseSchema{
		ResponseSchema: ResponseMessage("Coins fetched successfully"),
	}.Init(coins)
	return c.Status(200).JSON(response)
}

// @Summary Buy Coins
// @Description This endpoint allows a user to buy coins
// @Tags Wallet
// @Param coin body schemas.BuyCoinSchema true "Payment object"
// @Success 200 {object} schemas.PaymentResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /wallet/coins [post]
// @Security BearerAuth
func (ep Endpoint) BuyCoins(c *fiber.Ctx) error {
	stripe.Key = cfg.StripeSecretKey
	db := ep.DB
	user := RequestUser(c)

	data := schemas.BuyCoinSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	coin := models.Coin{}
	db.Where("id = ?", data.CoinID).Take(&coin)
	if coin.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No set of coins with that ID"))
	}

	var transaction models.Transaction
	if data.PaymentType == choices.PTYPE_STRIPE {
		// Create payment intent
		trans, errD := CreatePaymentIntent(db, *user, coin)
		if errD != nil {
			return c.Status(500).JSON(errD)
		}
		transaction = *trans
	}

	response := schemas.PaymentResponseSchema{
		ResponseSchema: ResponseMessage("Payment Data Generated"),
		Data:           schemas.TransactionSchema{}.Init(transaction),
	}
	return c.Status(200).JSON(response)
}
