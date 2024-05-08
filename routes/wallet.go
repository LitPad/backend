package routes

import (
	"encoding/json"
	"log"

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

// @Summary View Current Transactions
// @Description This endpoint returns all transactions of a user
// @Tags Wallet
// @Param page query int false "Current Page" default(1)
// @Param payment_status query string false "Payment Status"
// @Success 200 {object} schemas.TransactionsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /wallet/transactions [get]
// @Security BearerAuth
func (ep Endpoint) AllUserTransactions(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	transactions := []models.Transaction{}
	filterData := models.Transaction{UserID: user.ID}

	paymentStatus, errD := ValidatePaymentStatus(c)
	if errD != nil {
		return c.Status(400).JSON(errD)
	}
	empty := ""
	if paymentStatus != &empty {
		filterData.PaymentStatus = choices.PaymentStatus(*paymentStatus)
	}
	db.Where(filterData).Order("created_at DESC").Joins("Coin").Find(&transactions)
	// Paginate and return transactions
	paginatedData, paginatedTransactions, err := PaginateQueryset(transactions, c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	transactions = paginatedTransactions.([]models.Transaction)
	response := schemas.TransactionsResponseSchema{
		ResponseSchema: ResponseMessage("Transactions fetched successfully"),
		Data: schemas.TransactionsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(transactions),
	}
	return c.Status(200).JSON(response)
}

func (ep Endpoint) VerifyPayment(c *fiber.Ctx) error {
	stripe.Key = cfg.StripeSecretKey
	db := ep.DB
	transaction := models.Transaction{}
	event := stripe.Event{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &event); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment was successful
		transaction.Reference = paymentIntent.ID
		db.Joins("User").Joins("Coin").Take(&transaction, transaction)
		if transaction.ID != uuid.Nil {
			user := transaction.User
			coin := transaction.Coin
			user.Coins = user.Coins + coin.Amount
			transaction.PaymentStatus = choices.PSSUCCEEDED
			db.Save(&user)
			db.Save(&transaction)
		}
	case "payment_intent.canceled":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment was canceled
		transaction.Reference = paymentIntent.ID
		db.Model(&transaction).Update("payment_status", choices.PSCANCELED)
	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment failed
		transaction.Reference = paymentIntent.ID
		db.Model(&transaction).Update("payment_status", choices.PSFAILED)
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}
	return c.SendStatus(fiber.StatusOK)
}
