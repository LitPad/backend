package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/LitPad/backend/models/choices"
)

type Coin struct {
	BaseModel
	Amount int            `gorm:"default:0" json:"amount"`
	Price  decimal.Decimal `gorm:"default:0" json:"price"`
}

type Transaction struct {
	BaseModel
	Reference string    `json:"reference" gorm:"type: varchar(255);not null"`
	UserID    uuid.UUID `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	CoinID   uuid.UUID `json:"coin_id"`
	Coin     Coin      `gorm:"foreignKey:CoinID;constraint:OnDelete:CASCADE"`

	Verified bool      `json:"verified" gorm:"default:false"`
	PaymentType choices.PaymentType `json:"payment_type"`
}

type BoughtBooks struct {
	BaseModel
	ReaderID    uuid.UUID 
	Reader      User      `gorm:"foreignKey:ReaderID;constraint:OnDelete:CASCADE"`

	BookID   uuid.UUID 
	Book     Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
}
