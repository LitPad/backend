package schemas

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

// REQUEST BODY SCHEMAS
type RegisterUser struct {
	Email    string `json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
	Password string `json:"password" validate:"required,min=8,max=50" example:"strongpassword"`
}

type EmailRequestSchema struct {
	Email string `json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
}

type VerifyEmailRequestSchema struct {
	EmailRequestSchema
	Otp uint `json:"otp" validate:"required" example:"123456"`
}

type SetNewPasswordSchema struct {
	EmailRequestSchema
	TokenString string `json:"token_string" validate:"required" example:"Z2ZBYWjwXGXtCin3QnnABCHVfys6bcGPH49GrJEMtFIDQcU9TVL1AURNItZoBcTowOOeQMHofbp6WTxpYPlucdUEImQNWzMtH0ll"`
	Password string `json:"password" validate:"required,min=8,max=50" example:"newstrongpassword"`
}

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@email.com"`
	Password string `json:"password" validate:"required" example:"password"`
}

type SocialLoginSchema struct {
	DeviceType  choices.DeviceType `json:"device_type" validate:"device_type_validator"` 
	Token string `json:"token" validate:"required,min=10" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

type RefreshTokenSchema struct {
	Refresh string `json:"refresh" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

// RESPONSE BODY SCHEMAS
type RegisterResponseSchema struct {
	ResponseSchema
	Data EmailRequestSchema `json:"data"`
}

type TokensResponseSchema struct {
	UserProfile
	Access  string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
	Refresh string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

func (t TokensResponseSchema) Init(user models.User) TokensResponseSchema {
	t.UserProfile = t.UserProfile.Init(user)
	return t
}

type LoginResponseSchema struct {
	ResponseSchema
	Data TokensResponseSchema `json:"data"`
}
