package models

import "github.com/google/uuid"

type Waitlist struct {
	BaseModel
	Name    string `gorm:"varchar(1000)"`
	Email   string `gorm:"varchar(10000)"`
	GenreID uuid.UUID
	Genre   Genre `gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE;<-:false"`
}
