package schemas

// REQUEST BODY SCHEMAS
type RegisterUser struct {
	FirstName      string `json:"first_name" validate:"required,max=50" example:"John"`
	LastName       string `json:"last_name" validate:"required,max=50" example:"Doe"`
	Username       string `json:"username" validate:"required,max=1000" example:"john-doe"`
	Email          string `json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
	Password       string `json:"password" validate:"required,min=8,max=50" example:"strongpassword"`
	TermsAgreement bool   `json:"terms_agreement" validate:"eq=true"`
}

type EmailRequestSchema struct {
	Email string `json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
}

type VerifyEmailRequestSchema struct {
	TokenString					string					`json:"token_string" validate:"required" example:"Z2ZBYWjwXGXtCin3QnnABCHVfys6bcGPH49GrJEMtFIDQcU9TVL1AURNItZoBcTowOOeQMHofbp6WTxpYPlucdUEImQNWzMtH0ll"`
}

type SetNewPasswordSchema struct {
	VerifyEmailRequestSchema
	Password			string				`json:"password" validate:"required,min=8,max=50" example:"newstrongpassword"`
}

type LoginSchema struct {
	Email				string				`json:"email" validate:"required,email" example:"johndoe@email.com"`
	Password			string				`json:"password" validate:"required" example:"password"`
}

type SocialLoginSchema struct {
	Token				string				`json:"token" validate:"required,min=10" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

type RefreshTokenSchema struct {
	Refresh			string					`json:"refresh" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

// RESPONSE BODY SCHEMAS
type RegisterResponseSchema struct {
	ResponseSchema
	Data EmailRequestSchema `json:"data"`
}

type TokensResponseSchema struct {
	Access			string					`json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
	Refresh			string					`json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

type LoginResponseSchema struct {
	ResponseSchema
	Data			TokensResponseSchema		`json:"data"`
}