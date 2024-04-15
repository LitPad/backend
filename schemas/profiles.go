package schemas

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

type UserProfile struct {
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    Username    string `json:"username"`
    Email       string `json:"email"`
    Avatar      *string `json:"avatar"`
    Bio         *string `json:"bio"`
    AccountType choices.AccType `json:"account_type"`
    Followers   []models.User `json:"followers"`
    Followings  []models.User `json:"followings"`
}

func (u UserProfile) Init(user models.User) UserProfile {
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Username = user.Username
	u.Email = user.Email
	u.Avatar = user.Avatar
	u.Bio = user.Bio
	u.AccountType = user.AccountType
	u.Followers = user.Followers
	u.Followings = user.Followings
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
