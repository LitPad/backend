package models

import (
	"time"

	"github.com/LitPad/backend/config"
	"github.com/google/uuid"
)

var cfg = config.GetConfig()

type BaseModel struct {
	ID 			uuid.UUID `json:"-" gorm:"type:uuid;primarykey;not null;default:uuid_generate_v4()"`
	CreatedAt 	time.Time `json:"-" gorm:"not null"`
	UpdatedAt 	time.Time `json:"-" gorm:"not null"`
}