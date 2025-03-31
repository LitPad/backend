package routes

import (
	"fmt"
	"log"
	"strings"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"gorm.io/gorm"
)

func ResponseMessage(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}

func RequestUser(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func Session(c *fiber.Ctx, store *session.Store) *session.Session {
	// Get session from storage
	sess, err := store.Get(c)
	if err != nil {
		log.Println("Error Getting Session: ", err)
	}
	return sess
}

func GetBaseReferer(c *fiber.Ctx) string {
	referer := c.Context().Referer()
	return string(referer[:])
}

func CreatePaymentIntent(db *gorm.DB, user models.User, plan *models.SubscriptionPlan, paymentToken *string, coin *models.Coin, quantity int) (*models.Transaction, *utils.ErrorResponse) {
	cfg := config.GetConfig()
	stripe.Key = cfg.StripeSecretKey
	var price int64
	if coin != nil {
		price = coin.Price.Mul(decimal.NewFromFloat(100)).IntPart()
	} else {
		price = plan.Amount.Mul(decimal.NewFromFloat(100)).IntPart()
	}

	// Base PaymentIntent parameters
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(price),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	// Determine payment method type based on token presence
	if paymentToken != nil {
		// Google Pay (Token provided)
		params.PaymentMethodData = &stripe.PaymentIntentPaymentMethodDataParams{
			Type: stripe.String("card"),
			Card: &stripe.PaymentMethodCardParams{
				Token: stripe.String(*paymentToken), // Tokenized card
			},
		}
		params.ConfirmationMethod = stripe.String(string(stripe.PaymentIntentConfirmationMethodManual))
		params.Confirm = stripe.Bool(true)
	} else {
		// Card, Cashapp, etc (No token provided)
		params.PaymentMethodTypes = stripe.StringSlice([]string{"card", "cashapp"})
	}

	// Create the Payment Intent
	intent, err := paymentintent.New(params)
	if err != nil {
		errD := utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to create Payment Intent")
		return nil, &errD
	}

	// Create Transaction Object
	transaction := models.Transaction{
		Reference:    intent.ID,
		UserID:       user.ID,
		Quantity:     quantity,
		ClientSecret: intent.ClientSecret,
	}
	if coin != nil {
		transaction.CoinID = &coin.ID
		transaction.PaymentType = choices.PTYPE_STRIPE
		transaction.PaymentPurpose = choices.PP_COINS
	} else {
		transaction.SubscriptionPlanID = &plan.ID
		transaction.PaymentType = choices.PTYPE_GPAY
		transaction.PaymentPurpose = choices.PP_SUB
	}
	db.Create(&transaction)
	transaction.SubscriptionPlan = plan
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
		errD := utils.InvalidParamErr("Invalid payment status")
		return nil, &errD
	}
	return &status, nil
}

func CheckTagStrings(db *gorm.DB, submittedList []string) ([]models.Tag, *string) {
	tags := []models.Tag{}
	db.Find(&tags)
	// Create a map for quick lookup of predefined strings
	predefinedMap := make(map[string]bool)
	for _, item := range tags {
		predefinedMap[item.Slug] = true
	}

	// Iterate over the submitted list and check for any missing strings
	missingStrings := []string{}
	for _, item := range submittedList {
		if !predefinedMap[item] {
			missingStrings = append(missingStrings, item)
		}
	}

	// Return a message based on the result
	if len(missingStrings) > 0 {
		missingTags := strings.Join(missingStrings, ", ")
		errMsg := fmt.Sprintf("The following are invalid tag slugs: %v", missingTags)
		return tags, &errMsg
	}
	tagsToReturn := []models.Tag{}
	db.Where("slug IN ?", submittedList).Find(&tagsToReturn)
	return tagsToReturn, nil
}

func ReadBook(db *gorm.DB, bookID uuid.UUID, user *models.User) models.BookRead {
	bookRead := models.BookRead{UserID: user.ID, BookID: bookID}
	db.FirstOrCreate(&bookRead, bookRead)
	return bookRead
}

func IsAmongUserType(target string) bool {
	switch target {
	case "ADMIN", string(choices.ACCTYPE_READER), string(choices.ACCTYPE_AUTHOR):
		return true
	}
	return false
}

func GetQueryValue(c *fiber.Ctx, key string) *string {
	value := c.Query(key, "")
	if value == "" {
		return nil
	}
	return &value
}
