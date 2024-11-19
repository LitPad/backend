package routes

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	var transaction models.Transaction
	if data.PaymentType == choices.PTYPE_STRIPE {
		// Create payment intent
		trans, errD := CreateCheckoutSession(c, db, *user, &coin, data.Quantity, nil)
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
	sig := c.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(c.BodyRaw(), sig, cfg.StripeWebhookSecret)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"err": err})
	}

	// Handle different event types
	switch event.Type {
	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment was successful
		transaction.Reference = checkoutSession.ID
		db.Joins("User").Joins("Coin").Joins("SubscriptionPlan").Take(&transaction, transaction)
		if transaction.ID != uuid.Nil {
			// For Coins
			user := transaction.User
			if transaction.Coin != nil {
				coinsTotal := transaction.CoinsTotal()
				user.Coins = user.Coins + *coinsTotal
			} else {
				// For subscription
				subExpiry := time.Now().AddDate(0, 1, 0)
				if transaction.SubscriptionPlan.SubType == choices.ST_ANNUAL {
					subExpiry = time.Now().AddDate(0, 12, 0)
				}
				user.SubscriptionExpiry = &subExpiry
			}
			transaction.PaymentStatus = choices.PSSUCCEEDED
			db.Save(&user)
			db.Save(&transaction)
		}
	case "checkout.session.expired":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// Payment was canceled
		transaction.Reference = checkoutSession.ID
		db.Model(&transaction).Update("payment_status", choices.PSCANCELED)
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (ws WalletService) GetOnChainBalance(c *fiber.Ctx) error {
	accountID := c.Query("accountID")

	if len(accountID) == 0 {
		return c.Status(400).JSON(fiber.Map{"err": errors.New("Provide a valid account id")})
	}

	balance, err := ws.WS.GetBalance(accountID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"err": errors.New("Failed to retrieve balance")})
	}

	return c.JSON(fiber.Map{"balance": balance})
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

	var transaction models.Transaction
	// Create payment intent
	trans, errD := CreateCheckoutSession(c, db, *user, nil, 1, &plan)
	if errD != nil {
		return c.Status(500).JSON(errD)
	}
	transaction = *trans

	response := schemas.PaymentResponseSchema{
		ResponseSchema: ResponseMessage("Payment Data Generated"),
		Data:           schemas.TransactionSchema{}.Init(transaction),
	}
	return c.Status(200).JSON(response)
}
