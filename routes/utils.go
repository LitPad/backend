package routes

import (
	"fmt"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
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

func CreateCheckoutSession(c *fiber.Ctx, db *gorm.DB, user models.User, coin models.Coin, quantity int64) (*models.Transaction, *utils.ErrorResponse) {
	baseUrl := GetBaseReferer(c)
	stripe.Key = cfg.StripeSecretKey
	price := coin.Price.Mul(decimal.NewFromFloat(100)).IntPart()
	productName := fmt.Sprintf("%s coins", fmt.Sprint(coin.Amount))
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(baseUrl + cfg.StripeCheckoutSuccessUrlPath),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: &productName,
					},
					TaxBehavior: stripe.String(string(stripe.PriceTaxBehaviorExclusive)),
					UnitAmount:  stripe.Int64(price),
				},
				Quantity: stripe.Int64(quantity),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		CustomerEmail: &user.Email,
	}
	s, err := session.New(params)
	if err != nil {
		errD := utils.RequestErr(utils.ERR_SERVER_ERROR, "Something went wrong")
		return nil, &errD
	}

	// Create Transaction Object
	transaction := models.Transaction{Reference: s.ID, UserID: user.ID, CoinID: coin.ID, PaymentType: choices.PTYPE_STRIPE, Quantity: quantity, CheckoutURL: s.URL}
	db.Create(&transaction)
	transaction.Coin = coin
	return &transaction, nil
}

func IsValidPaymentStatus(s string) bool {
	switch choices.PaymentStatus(s) {
	case choices.PSPENDING, choices.PSSUCCEEDED, choices.PSFAILED, choices.PSCANCELED:
		return true
	}
	return false
}

func ValidatePaymentStatus(c *fiber.Ctx) (*string, *utils.ErrorResponse) {
	status := c.Query("payment_status", "")
	if status != "" && !IsValidPaymentStatus(status) {
		errD := utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid payment status")
		return nil, &errD
	}
	return &status, nil
}
