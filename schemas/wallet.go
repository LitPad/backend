package schemas

import (
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
	PaymentType choices.PaymentType `json:"payment_type" validate:"required,payment_type_validator" example:"STRIPE"`
	CoinID      uuid.UUID           `json:"coin_id" validate:"required" example:"19e8bd22-fab1-4bb4-ba82-77c41bea6b99"`
}

type TransactionSchema struct {
	Reference     string                `json:"reference"`
	ClientSecret  string                `json:"client_secret"`
	Coins         int                   `json:"coins" example:"10"`
	PaymentType   choices.PaymentType   `json:"payment_type" example:"STRIPR"`
	Amount        decimal.Decimal       `json:"amount" example:"10.35"`
	PaymentStatus choices.PaymentStatus `json:"payment_status"`
}

func (t TransactionSchema) Init(transaction models.Transaction) TransactionSchema {
	t.Reference = transaction.Reference
	t.ClientSecret = transaction.ClientSecret
	t.Coins = transaction.Coin.Amount
	t.PaymentType = transaction.PaymentType
	t.Amount = transaction.Coin.Price
	t.PaymentStatus = transaction.PaymentStatus
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
