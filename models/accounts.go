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
	FirstName         string          `json:"first_name" gorm:"type: varchar(255);not null" validate:"required,max=255" example:"John"`
	LastName          string          `json:"last_name" gorm:"type: varchar(255);not null" validate:"required,max=255" example:"Doe"`
	Username          string          `json:"username" gorm:"type: varchar(1000);not null;unique;" validate:"required,max=255" example:"john-doe"`
	Email             string          `json:"email" gorm:"not null;unique;" validate:"required,min=5,email" example:"johndoe@email.com"`
	Password          string          `json:"password" gorm:"not null" validate:"required,min=8,max=50" example:"strongpassword"`
	IsEmailVerified   bool            `json:"is_email_verified" gorm:"default:false" swaggerignore:"true"`
	IsSuperuser       bool            `json:"is_superuser" gorm:"default:false" swaggerignore:"true"`
	IsStaff           bool            `json:"is_staff" gorm:"default:false" swaggerignore:"true"`
	TermsAgreement    bool            `json:"terms_agreement" gorm:"default:false" validate:"eq=true"`
	Avatar            string          `gorm:"type:varchar(1000);null;" json:"avatar"`
	Access            *string         `gorm:"type:varchar(1000);null;" json:"access"`
	Refresh           *string         `gorm:"type:varchar(1000);null;" json:"refresh"`
	SocialLogin       bool            `gorm:"default:false"`
	Bio               *string         `gorm:"type:varchar(1000);null;" json:"bio"`
	AccountType       choices.AccType `gorm:"type:varchar(100); default:READER" json:"account_type"`
	Followings        []User          `gorm:"many2many:user_followers;foreignKey:ID;joinForeignKey:Follower;References:ID;joinReferences:Following"`
	Followers         []User          `gorm:"many2many:user_followers;foreignKey:ID;joinForeignKey:Following;References:ID;joinReferences:Follower"`
	Coins             int             `json:"coins" gorm:"default:0"`
	Lanterns          int             `json:"lanterns" gorm:"default:0"`
	LikeNotification  bool            `gorm:"default:false"`
	ReplyNotification bool            `gorm:"default:false"`

	// Back referenced
	Books []Book `gorm:"foreignKey:AuthorID"`
}

func (user User) BooksCount() int {
	return len(user.Books)
}

func (user User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Password = utils.HashPassword(user.Password)
	return
}

type Token struct {
	BaseModel
	UserId      uuid.UUID `json:"user_id" gorm:"unique"`
	User        User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	TokenString string    `json:"token_string"`
}

func (token *Token) BeforeSave(tx *gorm.DB) (err error) {
	token.TokenString = token.GenerateRandomToken(tx, "")
	return
}

func (token Token) GenerateRandomToken(db *gorm.DB, tokenString string) string {
	// Create new
	tokenStr := fmt.Sprintf("%s%s", utils.GetRandomString(100), tokenString)
	tokenData := Token{TokenString: tokenStr}
	if tokenData.ID != uuid.Nil {
		tokenStr = fmt.Sprintf("%s%s", tokenStr, utils.GetRandomString(6))
		return token.GenerateRandomToken(db, tokenStr)
	}
	return tokenStr
}

func (obj Token) CheckExpiration() bool {
	cfg := config.GetConfig()
	currentTime := time.Now().UTC()
	diff := int64(currentTime.Sub(obj.UpdatedAt).Seconds())
	emailExpirySecondsTimeout := cfg.EmailOtpExpireSeconds
	return diff > emailExpirySecondsTimeout
}

type Notification struct {
	BaseModel
	SenderID   uuid.UUID
	Sender     User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	ReceiverID uuid.UUID
	Receiver   User `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE"`
	Ntype      choices.NotificationTypeChoice
	Text       string

	BookID *uuid.UUID
	Book   *Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`

	ReviewID *uuid.UUID
	Review   *Review `gorm:"foreignKey:ReviewID;constraint:OnDelete:SET NULL;<-:false"`

	ReplyID *uuid.UUID
	Reply   *Reply `gorm:"foreignKey:ReplyID;constraint:OnDelete:SET NULL;<-:false"`

	SentGiftID *uuid.UUID
	SentGift   *SentGift `gorm:"foreignKey:SentGiftID;constraint:OnDelete:CASCADE;<-:false"`

	IsRead bool `gorm:"default:false"`
}
