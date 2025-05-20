package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/senders"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Register a new user
// @Description `This endpoint registers new users into our application.`
// @Tags Auth
// @Param user body schemas.RegisterUser true "User data"
// @Success 201 {object} schemas.RegisterResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (ep Endpoint) Register(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.RegisterUser{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	existingUser := models.User{}
	// Validate email uniqueness
	db.Take(&existingUser, models.User{Email: data.Email})
	if existingUser.ID != uuid.Nil {
		return c.Status(422).JSON(utils.ValidationErr("email", "Email already taken!"))
	}

	user := models.User{Email: data.Email, Password: data.Password}

	// Create User
	db.Save(&user)
	// Send Email
	go senders.SendEmail(&user, senders.ET_ACTIVATE, user.Otp, nil, nil)

	response := schemas.RegisterResponseSchema{
		ResponseSchema: ResponseMessage("Registration successful"),
		Data:           schemas.EmailRequestSchema{Email: user.Email},
	}
	return c.Status(201).JSON(response)
}

// @Summary Verify a user's email
// @Description `This endpoint verifies a user's email.`
// @Tags Auth
// @Param verify_email body schemas.VerifyEmailRequestSchema true "Verify Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/verify-email [post]
func (ep Endpoint) VerifyEmail(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.VerifyEmailRequestSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Email: data.Email, Otp: &data.Otp}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_OTP, "Invalid Email or OTP"))
	}

	if user.IsOtpExpired() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_EXPIRED_OTP, "Expired OTP"))
	}
	// Update User
	if user.IsEmailVerified {
		return c.Status(200).JSON(ResponseMessage("Email already verified"))
	}
	user.IsEmailVerified = true
	user.Otp = nil
	user.OtpExpiry = nil
	db.Save(&user)

	// Send Welcome Email
	go senders.SendEmail(&user, senders.ET_WELCOME, nil, nil, nil)
	return c.Status(200).JSON(ResponseMessage("Account verification successful"))
}

// @Summary Resend Verification Email
// @Description `This endpoint resends new otp to the user's email.`
// @Tags Auth
// @Param email body schemas.EmailRequestSchema true "Email data"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/resend-verification-email [post]
func (ep Endpoint) ResendVerificationEmail(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.EmailRequestSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Email: data.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_EMAIL, "Incorrect Email"))
	}

	if user.IsEmailVerified {
		return c.Status(200).JSON(schemas.ResponseSchema{Message: "Email already verified"}.Init())
	}

	// Send Email
	user.GenerateOTP(db)
	db.Save(&user)
	go senders.SendEmail(&user, senders.ET_ACTIVATE, user.Otp, nil, nil)
	return c.Status(200).JSON(ResponseMessage("Verification email sent"))
}

// @Summary Send Password Reset Link
// @Description `This endpoint sends new password reset link to the user's email.`
// @Tags Auth
// @Param email body schemas.EmailRequestSchema true "Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/send-password-reset-link [post]
func (ep Endpoint) SendPasswordResetLink(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.EmailRequestSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Email: data.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_EMAIL, "Incorrect Email"))
	}

	// Send Email
	user.GenerateToken(db)
	db.Save(&user)
	go senders.SendEmail(&user, senders.ET_RESET, nil, user.TokenString, nil)
	return c.Status(200).JSON(ResponseMessage("Password reset link sent"))
}

// @Summary Check Password Reset Token Validity
// @Description `This endpoint checks the validity of a password reset token.`
// @Tags Auth
// @Param token_string path string true "Token string"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/verify-password-reset-token/{token_string} [get]
func (ep Endpoint) VerifyPasswordResetToken(c *fiber.Ctx) error {
	db := ep.DB

	tokenStr := c.Params("token_string")

	user := models.User{TokenString: &tokenStr}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_TOKEN, "Invalid Token"))
	}

	if user.IsTokenExpired() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_EXPIRED_TOKEN, "Expired Token"))
	}
	return c.Status(200).JSON(ResponseMessage("Token verified successfully"))
}

// @Summary Set New Password
// @Description `This endpoint verifies the password reset token and set new password.`
// @Tags Auth
// @Param email body schemas.SetNewPasswordSchema true "Password reset object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/set-new-password [post]
func (ep Endpoint) SetNewPassword(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.SetNewPasswordSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Email: data.Email, TokenString: &data.TokenString}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_TOKEN, "Invalid Email or Token"))
	}

	if user.IsTokenExpired() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_EXPIRED_TOKEN, "Expired Token"))
	}

	// Set Password
	user.Password = utils.HashPassword(data.Password)
	user.TokenString = nil
	user.TokenExpiry = nil
	db.Save(&user)

	// Send Email
	go senders.SendEmail(&user, senders.ET_RESET_SUCC, nil, nil, nil)
	return c.Status(200).JSON(ResponseMessage("Password reset successful"))
}

// @Summary Login a user
// @Description This endpoint generates new access and refresh tokens for authentication
// @Tags Auth
// @Param user body schemas.LoginSchema true "User login"
// @Success 201 {object} schemas.LoginResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (ep Endpoint) Login(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.LoginSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Email: data.Email}
	db.Scopes(scopes.FollowerFollowingUnVerifiedPreloaderScope).Take(&user, user)
	if user.ID == uuid.Nil || !utils.CheckPasswordHash(data.Password, user.Password) {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_CREDENTIALS, "Invalid Credentials"))
	}

	if !user.IsEmailVerified {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNVERIFIED_USER, "Verify your email first"))
	}

	// Create Auth Tokens
	tokens := userManager.GenerateAuthTokens(db, user, GenerateAccessToken(user), GenerateRefreshToken())
	featuredContents := userManager.GetFeaturedContents(db, user)
	response := schemas.LoginResponseSchema{
		ResponseSchema: ResponseMessage("Login successful"),
		Data:           schemas.TokensResponseSchema{Access: tokens.Access, Refresh: tokens.Refresh}.Init(user, featuredContents),
	}
	return c.Status(201).JSON(response)
}

// @Summary Login a user via google
// @Description `This endpoint generates new access and refresh tokens for authentication via google`
// @Description `Pass in token gotten from gsi client authentication here in payload to retrieve tokens for authorization`
// @Tags Auth
// @Param user body schemas.SocialLoginSchema true "User login"
// @Success 201 {object} schemas.LoginResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/google [post]
func (ep Endpoint) GoogleLogin(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.SocialLoginSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	userGoogleData, errData := ConvertGoogleToken(data.Token, data.DeviceType)
	if errData != nil {
		return c.Status(401).JSON(errData)
	}

	email := userGoogleData.Email
	name := userGoogleData.Name
	avatar := userGoogleData.Picture

	user, token, err := RegisterSocialUser(db, email, name, &avatar)
	if err != nil {
		return c.Status(401).JSON(err)
	}
	featuredContents := userManager.GetFeaturedContents(db, *user)
	response := schemas.LoginResponseSchema{
		ResponseSchema: ResponseMessage("Social Login successful"),
		Data:           schemas.TokensResponseSchema{Access: token.Access, Refresh: token.Refresh}.Init(*user, featuredContents),
	}
	return c.Status(201).JSON(response)
}

// @Summary Login a user via facebook
// @Description `This endpoint generates new access and refresh tokens for authentication via facebook`
// @Description `Pass in token gotten from facebook client authentication here in payload to retrieve tokens for authorization`
// @Tags Auth
// @Param user body schemas.SocialLoginSchema true "User login"
// @Success 201 {object} schemas.LoginResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/facebook [post]
func (ep Endpoint) FacebookLogin(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.SocialLoginSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	userFacebookData, errData := ConvertFacebookToken(data.Token)
	if errData != nil {
		return c.Status(401).JSON(errData)
	}

	email := userFacebookData.Email
	name := userFacebookData.Name

	user, token, err := RegisterSocialUser(db, email, name, nil)
	if err != nil {
		return c.Status(401).JSON(err)
	}
	featuredContents := userManager.GetFeaturedContents(db, *user)
	response := schemas.LoginResponseSchema{
		ResponseSchema: ResponseMessage("Login successful"),
		Data:           schemas.TokensResponseSchema{Access: token.Access, Refresh: token.Refresh}.Init(*user, featuredContents),
	}
	return c.Status(201).JSON(response)
}

// @Summary Refresh tokens
// @Description This endpoint refresh tokens by generating new access and refresh tokens for a user
// @Tags Auth
// @Param refresh body schemas.RefreshTokenSchema true "Refresh token"
// @Success 201 {object} schemas.LoginResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/refresh [post]
func (ep Endpoint) Refresh(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.RefreshTokenSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	refreshToken := data.Refresh
	token := models.AuthToken{Refresh: refreshToken}
	db.Take(&token, token)
	if token.ID == uuid.Nil || !DecodeRefreshToken(refreshToken) {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, "Refresh token is invalid or expired"))
	}
	user := models.User{BaseModel: models.BaseModel{ID: token.UserID}}
	db.Find(scopes.FollowerFollowingPreloaderScope).Take(&user, user)

	// Create and Update Auth Tokens
	token.Access = GenerateAccessToken(user)
	token.Refresh = GenerateRefreshToken()
	db.Save(&token)
	featuredContents := userManager.GetFeaturedContents(db, user)

	response := schemas.LoginResponseSchema{
		ResponseSchema: ResponseMessage("Tokens refresh successful"),
		Data:           schemas.TokensResponseSchema{Access: token.Access, Refresh: token.Refresh}.Init(user, featuredContents),
	}
	return c.Status(201).JSON(response)
}

// @Summary Logout a user
// @Description This endpoint logs a user out from our application
// @Tags Auth
// @Success 200 {object} schemas.ResponseSchema
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/logout [get]
// @Security BearerAuth
func (ep Endpoint) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")[7:]
	db := ep.DB
	userManager.DeleteToken(db, token)
	return c.Status(200).JSON(ResponseMessage("Logout successful"))
}


// @Summary Logout a user from all devices
// @Description This endpoint logs a user out from every device
// @Tags Auth
// @Success 200 {object} schemas.ResponseSchema
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/logout/all [get]
// @Security BearerAuth
func (ep Endpoint) LogoutAll(c *fiber.Ctx) error {
	user := RequestUser(c)
	db := ep.DB
	userManager.DeleteAllToken(db, *user)
	return c.Status(200).JSON(ResponseMessage("Logout from all devices successful"))
}
