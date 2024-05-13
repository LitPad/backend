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
	Reference    string    `gorm:"type: varchar(1000);not null"` // payment id
	UserID       uuid.UUID `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	CoinID   uuid.UUID `json:"coin_id"`
	Coin     Coin      `gorm:"foreignKey:CoinID;constraint:OnDelete:CASCADE"`
	Quantity int64     `gorm:"default:1"`

	PaymentType   choices.PaymentType   `json:"payment_type"`
	PaymentStatus choices.PaymentStatus `json:"payment_status" gorm:"default:PENDING"`
	CheckoutURL		string
}

func (t Transaction) CoinsTotal() int {
	return t.Coin.Amount * int(t.Quantity)
}

type BoughtBooks struct {
	BaseModel
	ReaderID uuid.UUID
	Reader   User `gorm:"foreignKey:ReaderID;constraint:OnDelete:CASCADE"`

	BookID uuid.UUID
	Book   Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
}
