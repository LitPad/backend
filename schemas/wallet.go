package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CoinSchema struct {
	ID     uuid.UUID       `json:"id" example:"19e8bd22-fab1-4bb4-ba82-77c41bea6b99"`
	Amount int             `json:"amount" example:"5"`
	Price  decimal.Decimal `json:"price" example:"10.45"`
}

func (c CoinSchema) Init(coin models.Coin) CoinSchema {
	c.ID = coin.ID
	c.Amount = coin.Amount
	c.Price = coin.Price
	return c
}

type BuyCoinSchema struct {
	// PaymentType choices.PaymentType `json:"payment_type" validate:"required,payment_type_validator" example:"STRIPE"` // This should be stripe by default
	Quantity int       `json:"quantity" validate:"required" example:"2"`
	CoinID   uuid.UUID `json:"coin_id" validate:"required" example:"19e8bd22-fab1-4bb4-ba82-77c41bea6b99"`
}

type TransactionSchema struct {
	Reference      string                 `json:"reference"`
	Coins          *int                   `json:"coins" example:"10"`
	CoinsTotal     *int                   `json:"coins_total" example:"30"`
	PaymentType    choices.PaymentType    `json:"payment_type" example:"STRIPE"`
	PaymentPurpose choices.PaymentPurpose `json:"payment_purpose" example:"SUBSCRIPTION"`
	Quantity       int                    `json:"quantity" example:"10"`
	Amount         decimal.Decimal        `json:"amount" example:"10.35"`
	AmountTotal    decimal.Decimal        `json:"amount_total" example:"30.35"`
	PaymentStatus  choices.PaymentStatus  `json:"payment_status"`
	ClientSecret   string                 `json:"client_secret"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func (t TransactionSchema) Init(transaction models.Transaction) TransactionSchema {
	t.Reference = transaction.Reference
	if transaction.Coin != nil {
		amount := transaction.Coin.Amount
		t.Coins = &amount
		coinsTotal := transaction.Quantity * amount
		t.CoinsTotal = &coinsTotal
	}
	t.PaymentType = transaction.PaymentType
	t.PaymentStatus = transaction.PaymentStatus
	t.PaymentPurpose = transaction.PaymentPurpose
	if t.PaymentPurpose == choices.PP_COINS {
		t.Amount = transaction.Coin.Price
	} else {
		t.Amount = transaction.SubscriptionPlan.Amount
	}
	t.AmountTotal = t.Amount.Mul(decimal.NewFromInt(int64(transaction.Quantity)))
	t.Quantity = transaction.Quantity
	t.ClientSecret = transaction.ClientSecret
	t.CreatedAt = transaction.CreatedAt
	t.UpdatedAt = transaction.UpdatedAt
	return t
}

// RESPONSE SCHEMAS
type CoinsResponseSchema struct {
	ResponseSchema
	Data []CoinSchema `json:"data"`
}

func (c CoinsResponseSchema) Init(coins []models.Coin) CoinsResponseSchema {
	coinsData := []CoinSchema{}
	for i := range coins {
		coinsData = append(coinsData, CoinSchema{}.Init(coins[i]))
	}
	c.Data = coinsData
	return c
}

type PaymentResponseSchema struct {
	ResponseSchema
	Data TransactionSchema `json:"data"`
}

type TransactionsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []TransactionSchema `json:"transactions"`
}

func (t TransactionsResponseDataSchema) Init(transactions []models.Transaction) TransactionsResponseDataSchema {
	// Set Initial Data
	transactionItems := []TransactionSchema{}
	for _, transaction := range transactions {
		transactionItems = append(transactionItems, TransactionSchema{}.Init(transaction))
	}
	t.Items = transactionItems
	return t
}

type TransactionsResponseSchema struct {
	ResponseSchema
	Data TransactionsResponseDataSchema `json:"data"`
}

type CreateICPWallet struct{
	Username       string `json:"username" validate:"required,max=1000" example:"john-doe"`
}

type SubscriptionPlanSchema struct {
	Amount  decimal.Decimal                `json:"amount" validate:"required"`
	SubType choices.SubscriptionTypeChoice `json:"subtype" validate:"required,subscription_type_validator"`
}

func (s SubscriptionPlanSchema) Init(subscriptionPlan models.SubscriptionPlan) SubscriptionPlanSchema {
	s.Amount = subscriptionPlan.Amount
	s.SubType = subscriptionPlan.SubType
	return s
}

type SubscriptionPlansResponseSchema struct {
	ResponseSchema
	Data []SubscriptionPlanSchema `json:"data"`
}

func (s SubscriptionPlansResponseSchema) Init(subscriptionPlans []models.SubscriptionPlan) SubscriptionPlansResponseSchema {
	// Set Initial Data
	subscriptionPlanItems := []SubscriptionPlanSchema{}
	for _, plan := range subscriptionPlans {
		subscriptionPlanItems = append(subscriptionPlanItems, SubscriptionPlanSchema{}.Init(plan))
	}
	s.Data = subscriptionPlanItems
	return s
}

type SubscriptionPlanResponseSchema struct {
	ResponseSchema
	Data SubscriptionPlanSchema `json:"data"`
}

type CreateSubscriptionSchema struct {
	SubType            choices.SubscriptionTypeChoice `json:"subtype" validate:"required,subscription_type_validator"`
	PaymentMethodToken string                         `json:"payment_method_token" validate:"required"`
}

