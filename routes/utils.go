package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"gorm.io/gorm"
)

func ResponseMessage(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}

func RequestUser(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func GetBaseReferer(c *fiber.Ctx) string {
	referer := c.Context().Referer()
	return string(referer[:])
}

func CreatePaymentIntent(db *gorm.DB, user models.User, coin models.Coin) (*models.Transaction, *utils.ErrorResponse) {
	stripe.Key = cfg.StripeSecretKey
	price := coin.Price.Mul(decimal.NewFromFloat(100)).IntPart()
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(price)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		errD := utils.RequestErr(utils.ERR_SERVER_ERROR, "Something went wrong")
		return nil, &errD
	}

	// Create Transaction Object
	transaction := models.Transaction{Reference: pi.ID, ClientSecret: pi.ClientSecret, UserID: user.ID, CoinID: coin.ID, PaymentType: choices.PTYPE_STRIPE}
	db.Create(&transaction)
	transaction.Coin = coin 
	return &transaction, nil
}
