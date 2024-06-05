package models

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Gift struct {
	BaseModel
	Name     string	`gorm:"unique"`
	Slug     string	`gorm:"unique"`
	Price    int
	Image    string
	Lanterns int
}

func (g *Gift) BeforeSave(tx *gorm.DB) (err error) {
	g.Slug = slug.Make(g.Name)
	return
}

type SentGift struct {
	BaseModel
	SenderID uuid.UUID
	Sender   User      `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;<-:false"`

	ReceiverID uuid.UUID
	Receiver   User      `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE;<-:false"`

	GiftID uuid.UUID
	Gift   Gift      `gorm:"foreignKey:GiftID;constraint:OnDelete:CASCADE;<-:false"`

	Claimed bool	`gorm:"default:false"`
}