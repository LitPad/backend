package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name            *string `gorm:"type: varchar(255);null"`
	Username        string `gorm:"type: varchar(1000);not null;unique;"`
	Email           string `gorm:"not null;unique;"`
	Password        string `gorm:"not null"`
	IsEmailVerified bool   `gorm:"default:false"`
	IsSuperuser     bool   `gorm:"default:false"`
	IsStaff         bool   `gorm:"default:false"`
	IsActive        bool   `gorm:"default:true"`

	Otp         *uint      `gorm:"null"`
	OtpExpiry   *time.Time `gorm:"null"`
	TokenString *string    `gorm:"null"`
	TokenExpiry *time.Time `gorm:"null"`

	Avatar            string          `gorm:"type:varchar(1000);null;"`
	Access            *string         `gorm:"type:varchar(1000);null;"`
	Refresh           *string         `gorm:"type:varchar(1000);null;"`
	SocialLogin       bool            `gorm:"default:false"`
	Bio               *string         `gorm:"type:varchar(1000);null;"`
	AccountType       choices.AccType `gorm:"type:varchar(100); default:READER"`
	Followings        []User          `gorm:"many2many:user_followers;foreignKey:ID;joinForeignKey:Follower;References:ID;joinReferences:Following"`
	Followers         []User          `gorm:"many2many:user_followers;foreignKey:ID;joinForeignKey:Following;References:ID;joinReferences:Follower"`
	Coins             int             `gorm:"default:0"`
	Lanterns          int             `gorm:"default:0"`
	LikeNotification  bool            `gorm:"default:false"`
	ReplyNotification bool            `gorm:"default:false"`

	CurrentPlan        *choices.SubscriptionTypeChoice `gorm:"null"`
	SubscriptionExpiry *time.Time                      `gorm:"index,null"`
	ReminderSent       bool                            `gorm:"default:false"`

	// Back referenced
	Books []Book `gorm:"foreignKey:AuthorID"`
}

func (u *User) GenerateOTP(db *gorm.DB) {
	cfg := config.GetConfig()
	// Create new otp
	otp := utils.GetRandomInt(6)
	u.Otp = &otp
	expiry := time.Now().UTC().Add(time.Duration(cfg.EmailOtpExpireSeconds) * time.Second)
	u.OtpExpiry = &expiry
}

func (u *User) GenerateToken(db *gorm.DB) {
	cfg := config.GetConfig()
	// Create new token
	tokenString := utils.GetRandomString(70)
	u.TokenString = &tokenString
	expiry := time.Now().UTC().Add(time.Duration(cfg.EmailOtpExpireSeconds) * time.Second)
	u.TokenExpiry = &expiry
}

func (u User) IsTokenExpired() bool {
	if u.TokenExpiry == nil { return true }
	return time.Now().UTC().After((*u.TokenExpiry).UTC())
}

func (u User) IsOtpExpired() bool {
	if u.OtpExpiry == nil { return true }
	return time.Now().UTC().After((*u.OtpExpiry).UTC())
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

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Username = user.GenerateUsername(tx, user.Email, nil)
	user.GenerateOTP(tx)
	user.Password = utils.HashPassword(user.Password)
	return
}

func (user *User) GenerateUsername(db *gorm.DB, email string, username *string) string {
	emailSubstr := strings.Split(email, "@")[0]

	uniqueUsername := slug.Make(emailSubstr)
	if username != nil {
		uniqueUsername = *username
	}

	// Check for uniqueness and adjust if necessary
	for {
		exisitngUser := User{Username: uniqueUsername}
		db.Take(&exisitngUser, exisitngUser)
		if exisitngUser.ID == uuid.Nil {
			// Username is unique
			break
		}
		// Append a random string to make it unique
		randomStr := strconv.FormatUint(uint64(utils.GetRandomInt(7)), 10)
		uniqueUsername = slug.Make(emailSubstr) + randomStr
	}
	return uniqueUsername
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
	Review   *Comment `gorm:"foreignKey:ReviewID;constraint:OnDelete:SET NULL;<-:false"`

	ReplyID *uuid.UUID
	Reply   *Reply `gorm:"foreignKey:ReplyID;constraint:OnDelete:SET NULL;<-:false"`

	SentGiftID *uuid.UUID
	SentGift   *SentGift `gorm:"foreignKey:SentGiftID;constraint:OnDelete:CASCADE;<-:false"`

	IsRead bool `gorm:"default:false"`
}
