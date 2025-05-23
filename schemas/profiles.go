package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/google/uuid"
)

type FollowerData struct {
	Name           *string         `json:"name"`
	Username       string          `json:"username"`
	AccountType    choices.AccType `json:"account_type"`
	Avatar         *string         `json:"avatar"`
	FollowersCount int             `json:"followers_count"`
	StoriesCount   int             `json:"stories_count"`
}

func (dto FollowerData) FromModel(user models.User) FollowerData {
	dto.Name = user.Name
	dto.Username = user.Username
	dto.Avatar = &user.Avatar
	dto.AccountType = user.AccountType
	dto.FollowersCount = user.FollowersCount()
	dto.StoriesCount = user.BooksCount()
	return dto
}

type UserProfile struct {
	Name         *string                         `json:"name"`
	Username     string                          `json:"username"`
	Email        string                          `json:"email"`
	Avatar       *string                         `json:"avatar"`
	Bio          *string                         `json:"bio"`
	AccountType  choices.AccType                 `json:"account_type"`
	StoriesCount int                             `json:"stories_count"`
	Followers    []FollowerData                  `json:"followers"`
	Followings   []FollowerData                  `json:"followings"`
	CreatedAt    time.Time                       `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
	CurrentPlan  *choices.SubscriptionTypeChoice `json:"current_plan"`
	IsFollowing  bool                            `json:"is_following"`
	IsActive     bool                            `json:"is_active"`
}

func (u UserProfile) Init(user models.User, currentUser *models.User) UserProfile {
	followers := []FollowerData{}
	followings := []FollowerData{}
	isFollowing := false
	for _, follower := range user.Followers {
		followerData := FollowerData{}.FromModel(follower)
		followers = append(followers, followerData)
		if currentUser != nil && follower.ID == currentUser.ID {
			isFollowing = true
		}
	}
	for _, following := range user.Followings {
		followingData := FollowerData{}.FromModel(following)
		followings = append(followings, followingData)
	}
	u = UserProfile{
		Name:         user.Name,
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       &user.Avatar,
		Bio:          user.Bio,
		AccountType:  user.AccountType,
		Followers:    followers,
		Followings:   followings,
		StoriesCount: user.BooksCount(),
		CreatedAt:    user.CreatedAt,
		CurrentPlan:  user.CurrentPlan,
		IsFollowing:  isFollowing,
		IsActive:     user.IsActive,
	}
	return u
}

type UserProfileResponseSchema struct {
	ResponseSchema
	Data UserProfile `json:"data"`
}

type UpdateUserProfileSchema struct {
	Name     *string `json:"name,omitempty" validate:"min=3,max=1000" example:"John Doe"`
	Bio      *string `json:"bio,omitempty" validate:"min=3,max=1000" example:"I'm here to read good books"`
	Username *string `json:"username,omitempty" validate:"min=3,max=1000" example:"johndoe"`
}

type UpdateUserRoleSchema struct {
	AccountType choices.AccType `json:"account_type" validate:"account_type_validator" example:"WRITER"`
}

type UpdatePasswordSchema struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=50" example:"oldpassword"`
	OldPassword string `json:"old_password" validate:"required,min=8,max=50" example:"newstrongpassword"`
}

// NOTIFICATIONS
type NotificationBookSchema struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	CoverImage string `json:"cover_image"`
}

type NotificationSchema struct {
	ID         uuid.UUID                      `json:"id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`
	Sender     UserDataSchema                 `json:"sender"`
	ReceiverID *uuid.UUID                     `json:"receiver_id,omitempty"`
	Ntype      choices.NotificationTypeChoice `json:"ntype"`
	Text       string                         `json:"text"`
	Book       *NotificationBookSchema        `json:"book"`                                                        // Bought book, vote, comment and reply
	CommentID  *uuid.UUID                     `json:"comment_id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`   // reviewed, reply, like
	SentGiftID *uuid.UUID                     `json:"sent_gift_id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"` // If someone sent you a gift
	IsRead     bool                           `json:"is_read"`
	CreatedAt  time.Time                      `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (n NotificationSchema) Init(notification models.Notification, showReceiver ...bool) NotificationSchema {
	n.ID = notification.ID
	n.Sender = n.Sender.Init(notification.Sender)
	n.Ntype = notification.Ntype
	n.Text = notification.Text
	if notification.Book != nil {
		n.Book = &NotificationBookSchema{
			Title:      notification.Book.Title,
			Slug:       notification.Book.Slug,
			CoverImage: notification.Book.CoverImage,
		}
	}
	n.CommentID = notification.CommentID
	n.SentGiftID = notification.SentGiftID
	n.IsRead = notification.IsRead
	n.CreatedAt = notification.CreatedAt
	if len(showReceiver) > 0 {
		n.ReceiverID = &notification.ReceiverID
	}
	return n
}

type NotificationsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []NotificationSchema `json:"notifications"`
}

func (n NotificationsResponseDataSchema) Init(notifications []models.Notification) NotificationsResponseDataSchema {
	// Set Initial Data
	notificationItems := make([]NotificationSchema, 0)
	for _, notification := range notifications {
		notificationItems = append(notificationItems, NotificationSchema{}.Init(notification))
	}
	n.Items = notificationItems
	return n
}

type NotificationsResponseSchema struct {
	ResponseSchema
	Data NotificationsResponseDataSchema `json:"data"`
}

type ReadNotificationSchema struct {
	MarkAllAsRead bool       `json:"mark_all_as_read" example:"false"`
	ID            *uuid.UUID `json:"id" validate:"required_if=MarkAllAsRead false,omitempty" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
}
