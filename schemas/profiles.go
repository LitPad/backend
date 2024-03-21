package schemas

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

type UserProfile struct {
	FirstName   string          `json:"first_name" validate:"required,max=50" example:"John"`
	LastName    string          `json:"last_name" validate:"required,max=50" example:"Doe"`
	Username    string          `json:"username" validate:"required,max=1000" example:"john-doe"`
	Email       string          `json:"email" validate:"required"`
	Avatar      *string         `json:"avatar"`
	Bio         *string         `json:"bio"`
	AccountType choices.AccType `json:"account_type"`
}

func (u UserProfile) Init(user models.User) UserProfile {
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Username = user.Username
	u.Email = user.Email
	u.Avatar = user.Avatar
	u.Bio = user.Bio
	u.AccountType = user.AccountType
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
