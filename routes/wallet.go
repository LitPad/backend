package routes

import (
	"encoding/json"
	"log"
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/senders"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
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

	// Create payment intent
	transaction, errD := CreatePaymentIntent(db, *user, nil, nil, &coin, data.Quantity)
	if errD != nil {
		return c.Status(500).JSON(errD)
	}

	response := schemas.PaymentResponseSchema{
		ResponseSchema: ResponseMessage("Payment Data Generated"),
		Data:           schemas.TransactionSchema{}.Init(*transaction),
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
	db.Where(filterData).Order("created_at DESC").Joins("Coin").Joins("SubscriptionPlan").Find(&transactions)
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
	sig := c.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(c.BodyRaw(), sig, cfg.StripeWebhookSecret)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"err": err})
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		var intent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &intent)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment was successful
		transaction.Reference = intent.ID
		db.Joins("User").Joins("Coin").Joins("SubscriptionPlan").Take(&transaction, transaction)
		if transaction.ID != uuid.Nil {
			user := transaction.User
			subPlan := transaction.SubscriptionPlan
			if subPlan != nil {
				// For subscription
				emailD := map[string]interface{}{"amount": subPlan.Amount}
				expectedAmount := subPlan.Amount.Mul(decimal.NewFromFloat(100)).IntPart() // Convert to cents
				if intent.AmountReceived < expectedAmount {
					transaction.PaymentStatus = choices.PSFAILED
					go senders.SendEmail(&transaction.User, "payment-failed", nil, nil, emailD)
				} else {
					subExpiry := time.Now().AddDate(0, 1, 0)
					if transaction.SubscriptionPlan.SubType == choices.ST_ANNUAL {
						subExpiry = time.Now().AddDate(0, 12, 0)
					}
					user.SubscriptionExpiry = &subExpiry
					transaction.PaymentStatus = choices.PSSUCCEEDED
					go senders.SendEmail(&transaction.User, "payment-succeeded", nil, nil, emailD)
				}
			} else {
				coin := transaction.Coin
				emailD := map[string]interface{}{"amount": coin.Price}
				expectedAmount := coin.Price.Mul(decimal.NewFromFloat(100)).IntPart() // Convert to cents
				if intent.AmountReceived < expectedAmount {
					transaction.PaymentStatus = choices.PSFAILED
					go senders.SendEmail(&transaction.User, "payment-failed", nil, nil, emailD)
				} else {
					coinsTotal := transaction.CoinsTotal()
					user.Coins = user.Coins + *coinsTotal
					transaction.PaymentStatus = choices.PSSUCCEEDED
					go senders.SendEmail(&transaction.User, "payment-succeeded", nil, nil, emailD)
				}
			}
			db.Save(&user)
			db.Save(&transaction)
		}
	case "payment_intent.payment_failed":
		var intent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &intent)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment failed
		transaction.Reference = intent.ID
		db.Joins("User").Joins("Coin").Joins("SubscriptionPlan").Take(&transaction, transaction)
		if transaction.ID != uuid.Nil {
			transaction.PaymentStatus = choices.PSFAILED
			db.Save(&transaction)
			coin := transaction.Coin
			plan := transaction.SubscriptionPlan
			var amount decimal.Decimal
			if coin != nil {
				amount = coin.Price
			} else {
				amount = plan.Amount
			}
			emailD := map[string]interface{}{"amount": amount}
			go senders.SendEmail(&transaction.User, "payment-failed", nil, nil, emailD)
		}
	case "payment_intent.canceled":
		var intent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &intent)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment canceled
		transaction.Reference = intent.ID
		db.Joins("User").Joins("Coin").Joins("SubscriptionPlan").Take(&transaction, transaction)
		if transaction.ID != uuid.Nil {
			transaction.PaymentStatus = choices.PSCANCELED
			db.Save(&transaction)
			coin := transaction.Coin
			plan := transaction.SubscriptionPlan
			var amount decimal.Decimal
			if coin != nil {
				amount = coin.Price
			} else {
				amount = plan.Amount
			}
			emailD := map[string]interface{}{"amount": amount}
			go senders.SendEmail(&transaction.User, "payment-canceled", nil, nil, emailD)
		}
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}
	return c.SendStatus(fiber.StatusOK)
}

// @Summary List Available Subscription Plans
// @Description Retrieves a list of available subscription plans.
// @Tags Wallet
// @Accept json
// @Produce json
// @Success 200 {object} schemas.SubscriptionPlansResponseSchema "Successfully retrieved list of plans"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /wallet/plans [get]
func (ep Endpoint) GetSubscriptionPlans(c *fiber.Ctx) error {
	db := ep.DB
	plans := []models.SubscriptionPlan{}
	db.Find(&plans)
	response := schemas.SubscriptionPlansResponseSchema{
		ResponseSchema: ResponseMessage("Plans fetched successfully"),
	}.Init(plans)
	return c.Status(200).JSON(response)
}

// @Summary Update A Plan Amount
// @Description This endpoint allows an admin to change the amount of a plan
// @Tags Wallet
// @Param plan body schemas.SubscriptionPlanSchema true "Plan data"
// @Success 200 {object} schemas.SubscriptionPlanResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /wallet/plans [put]
// @Security BearerAuth
func (ep Endpoint) UpdateSubscriptionPlan(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.SubscriptionPlanSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	plan := models.SubscriptionPlan{SubType: data.SubType}
	db.Take(&plan, plan)
	plan.Amount = data.Amount
	db.Save(&plan)
	response := schemas.SubscriptionPlanResponseSchema{
		ResponseSchema: ResponseMessage("Plan updated successfully"),
		Data:           schemas.SubscriptionPlanSchema{}.Init(plan),
	}
	return c.Status(200).JSON(response)
}

// @Summary Subscribe
// @Description This endpoint allows a user to create a subscription for books
// @Tags Wallet
// @Param subscription body schemas.CreateSubscriptionSchema true "Payment object"
// @Success 200 {object} schemas.PaymentResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /wallet/subscription [post]
// @Security BearerAuth
func (ep Endpoint) BookSubscription(c *fiber.Ctx) error {
	stripe.Key = cfg.StripeSecretKey
	db := ep.DB
	user := RequestUser(c)

	data := schemas.CreateSubscriptionSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	plan := models.SubscriptionPlan{}
	db.Where("type = ?", data.SubType).Take(&plan)
	if plan.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No subscription plan with that type"))
	}

	// Create payment intent
	transaction, errD := CreatePaymentIntent(db, *user, &plan, &data.PaymentMethodToken, nil, 1)
	if errD != nil {
		return c.Status(500).JSON(errD)
	}

	response := schemas.PaymentResponseSchema{
		ResponseSchema: ResponseMessage("Payment Data Generated"),
		Data:           schemas.TransactionSchema{}.Init(*transaction),
	}
	return c.Status(200).JSON(response)
}
