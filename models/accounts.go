package models

import (
	"fmt"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	FirstName       string          `json:"first_name" gorm:"type: varchar(255);not null" validate:"required,max=255" example:"John"`
	LastName        string          `json:"last_name" gorm:"type: varchar(255);not null" validate:"required,max=255" example:"Doe"`
	Username        string          `json:"username" gorm:"type: varchar(1000);not null;unique;" validate:"required,max=255" example:"john-doe"`
	Email           string          `json:"email" gorm:"not null;unique;" validate:"required,min=5,email" example:"johndoe@email.com"`
	Password        string          `json:"password" gorm:"not null" validate:"required,min=8,max=50" example:"strongpassword"`
	IsEmailVerified bool            `json:"is_email_verified" gorm:"default:false" swaggerignore:"true"`
	IsSuperuser     bool            `json:"is_superuser" gorm:"default:false" swaggerignore:"true"`
	IsStaff         bool            `json:"is_staff" gorm:"default:false" swaggerignore:"true"`
	TermsAgreement  bool            `json:"terms_agreement" gorm:"default:false" validate:"eq=true"`
	Avatar          *string         `gorm:"type:varchar(1000);null;" json:"avatar"`
	Access          *string         `gorm:"type:varchar(1000);null;" json:"access"`
	Refresh         *string         `gorm:"type:varchar(1000);null;" json:"refresh"`
	SocialLogin     bool            `gorm:"default:false"`
	Bio             *string         `gorm:"type:varchar(1000);null;" json:"bio"`
	AccountType     choices.AccType `gorm:"type:varchar(100); default:READER" json:"account_type"`
	Followers		[]User			`json:"followers" gorm:"many2many:user_followers;"`
	Followings		[]User			`json:"followings" gorm:"many2many:user_followings;"`
}

func (user User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Password = utils.HashPassword(user.Password)
	return
}

type Otp struct {
	BaseModel
	UserId uuid.UUID `json:"user_id" gorm:"unique"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Code   uint32    `json:"code"`
}

func (otp *Otp) BeforeSave(tx *gorm.DB) (err error) {
	code := uint32(utils.GetRandomInt(6))
	otp.Code = code
	return
}

func (obj Otp) CheckExpiration() bool {
	cfg := config.GetConfig()
	currentTime := time.Now().UTC()
	diff := int64(currentTime.Sub(obj.UpdatedAt).Seconds())
	emailExpirySecondsTimeout := cfg.EmailOtpExpireSeconds
	return diff > emailExpirySecondsTimeout
}

