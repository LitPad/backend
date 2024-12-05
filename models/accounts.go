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
	FirstName          string          `gorm:"type: varchar(255);not null"`
	LastName           string          `gorm:"type: varchar(255);not null"`
	Username           string          `gorm:"type: varchar(1000);not null;unique;"`
	Email              string          `gorm:"not null;unique;"`
	Password           string          `gorm:"not null"`
	IsEmailVerified    bool            `gorm:"default:false"`
	IsSuperuser        bool            `gorm:"default:false"`
	IsStaff            bool            `gorm:"default:false"`
	IsActive           bool            `gorm:"default:true"`
	TermsAgreement     bool            `gorm:"default:false"`
	Avatar             string          `gorm:"type:varchar(1000);null;"`
	Access             *string         `gorm:"type:varchar(1000);null;"`
	Refresh            *string         `gorm:"type:varchar(1000);null;"`
	SocialLogin        bool            `gorm:"default:false"`
	Bio                *string         `gorm:"type:varchar(1000);null;"`
	AccountType        choices.AccType `gorm:"type:varchar(100); default:READER"`
	Followings         []User          `gorm:"many2many:user_followers;foreignKey:ID;joinForeignKey:Follower;References:ID;joinReferences:Following"`
	Followers          []User          `gorm:"many2many:user_followers;foreignKey:ID;joinForeignKey:Following;References:ID;joinReferences:Follower"`
	Coins              int             `gorm:"default:0"`
	Lanterns           int             `gorm:"default:0"`
	LikeNotification   bool            `gorm:"default:false"`
	ReplyNotification  bool            `gorm:"default:false"`
	SubscriptionExpiry *time.Time      `gorm:"null"`

	// Back referenced
	Books []Book `gorm:"foreignKey:AuthorID"`
}

func (user User) SubscriptionExpired() bool {
	if user.SubscriptionExpiry == nil {
		return true
	}
	return time.Now().After(*user.SubscriptionExpiry)
}

func (user User) BooksCount() int {
	return len(user.Books)
}

func (user User) FollowersCount() int {
	return len(user.Followers)
}

func (user User) FollowingsCount() int {
	return len(user.Followings)
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
	UserId      uuid.UUID `gorm:"unique"`
	User        User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	TokenString string
}

func (token *Token) BeforeSave(tx *gorm.DB) (err error) {
	token.TokenString = token.GenerateRandomToken(tx)
	return
}

func (token Token) GenerateRandomToken(db *gorm.DB) string {
	// Create new
	tokenStr := utils.GetRandomString(100)
	tokenData := Token{TokenString: tokenStr}
	db.Take(&tokenData, tokenData)
	if tokenData.ID != uuid.Nil {
		return token.GenerateRandomToken(db)
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
	Sender     User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;<-:false"`
	ReceiverID uuid.UUID
	Receiver   User `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE;<-:false"`
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
