package schemas

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

type FollowerData struct {
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	AccountType choices.AccType `json:"account_type"`
	Avatar         *string `json:"avatar"`
	FollowersCount int     `json:"followers_count"`
	StoriesCount   int     `json:"stories_count"`
}

func (dto FollowerData) FromModel(user models.User) FollowerData {
	dto.Name = user.FullName()
	dto.Username = user.Username
	dto.Avatar = user.Avatar
	dto.AccountType = user.AccountType
	dto.FollowersCount = len(user.Followers)
	dto.StoriesCount = len(user.Books)
	return dto
}

type UserProfile struct {
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	Avatar      *string         `json:"avatar"`
	Bio         *string         `json:"bio"`
	AccountType choices.AccType `json:"account_type"`
	Followers   []FollowerData  `json:"followers"`
	Followings  []FollowerData  `json:"followings"`
}

func (u UserProfile) Init(user models.User) UserProfile {
	followers := []FollowerData{}
	followings := []FollowerData{}
	for _, follower := range user.Followers {
		followerData := FollowerData{}.FromModel(follower)
		followers = append(followers, followerData)
	}

	for _, following := range user.Followings {
		followingData := FollowerData{}.FromModel(following)
		followings = append(followings, followingData)
	}

	u = UserProfile{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		AccountType: user.AccountType,
		Followers:   followers,
		Followings:  followings,
	}
	return u
}

type UserProfileResponseSchema struct {
	ResponseSchema
	Data UserProfile `json:"data"`
}

type UpdateUserProfileSchema struct {
	// Bio				*string `json:"bio"`
	Username *string `json:"username,omitempty" validate:"min=3,max=1000" example:"john-doe"`
}

type UpdatePasswordSchema struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=50" example:"oldpassword"`
	OldPassword string `json:"old_password" validate:"required,min=8,max=50" example:"newstrongpassword"`
}
