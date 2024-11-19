package models

import (
	"github.com/LitPad/backend/models/choices"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Coin struct {
	BaseModel
	Amount int             `gorm:"default:0" json:"amount"`
	Price  decimal.Decimal `gorm:"default:0" json:"price"`
}

type Transaction struct {
	BaseModel
	Reference string    `gorm:"type: varchar(1000);not null"` // payment id
	UserID    uuid.UUID `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	// FOR COINS
	CoinID   *uuid.UUID `json:"coin_id"`
	Coin     *Coin      `gorm:"foreignKey:CoinID;constraint:OnDelete:SET NULL"`
	Quantity int        `gorm:"default:1"`
	// -----------------

	// FOR SUBSCRIPTION
	SubscriptionPlanID *uuid.UUID        `json:"subscription_plan_id"`
	SubscriptionPlan   *SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID;constraint:OnDelete:SET NULL"`

	// ---------------------

	PaymentType    choices.PaymentType    `json:"payment_type"`
	PaymentPurpose choices.PaymentPurpose `json:"payment_purpose"`
	PaymentStatus  choices.PaymentStatus  `json:"payment_status" gorm:"default:PENDING"`
	CheckoutURL    string
}

func (t Transaction) CoinsTotal() *int {
	if t.Coin != nil {
		amount := t.Coin.Amount
		coinsTotal := t.Quantity * amount
		return &coinsTotal
	}
	return nil
}

type SubscriptionPlan struct {
	BaseModel
	Amount  decimal.Decimal                `gorm:"default:0"`
	SubType choices.SubscriptionTypeChoice `gorm:"default:MONTHLY;unique"`
}
