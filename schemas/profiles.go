package schemas

type AccType string


type UserProfile struct{
	FirstName      string `json:"first_name" validate:"required,max=50" example:"John"`
	LastName       string `json:"last_name" validate:"required,max=50" example:"Doe"`
	Username       string `json:"username" validate:"required,max=1000" example:"john-doe"`
	Email          string `json:"email" validate:"required"`
	Avatar          *string `json:"avatar"`
	Bio				*string `json:"bio"`
	AccountType     AccType `json:"account_type"`
}

type UserProfileResponseSchema struct{
	ResponseSchema
	Data UserProfile
}

type UpdateUserProfileSchema struct{
	// Bio				*string `json:"bio"`
	Username       string `json:"username" validate:"max=1000" example:"john-doe"`
}

type UpdatePasswordSchema struct{
	NewPassword			string				`json:"newPassword" validate:"required,min=8,max=50" example:"newstrongpassword"`
	OldPassword			string				`json:"oldPassword" validate:"required,min=8,max=50" example:"newstrongpassword"`
}