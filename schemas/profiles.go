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
    Followers   []UserProfile `json:"followers"`
    Followings  []UserProfile `json:"followings"`
}

func (dto *UserProfile) FromModel(user models.User) {
    dto.FirstName = user.FirstName
    dto.LastName = user.LastName
    dto.Username = user.Username
    dto.Avatar = user.Avatar
    dto.Bio = user.Bio
    dto.AccountType = user.AccountType
}

func (u UserProfile) Init (user models.User) UserProfile{
	var followers []UserProfile
	var followings []UserProfile

	for _, follower := range user.Followers{
		var dto UserProfile
		dto.FromModel(follower)
		followers = append(followers, dto)
	}

	 for _, following := range user.Followings {
        var dto UserProfile
		dto.FromModel(following)
        followings = append(followings, dto)
    }

	u = UserProfile{
		FirstName : user.FirstName,
		LastName : user.LastName,
		Username : user.Username,
		Email : user.Email,
		Avatar : user.Avatar,
		Bio : user.Bio,
		AccountType : user.AccountType,
		Followers : followers,
		Followings : followings,
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
