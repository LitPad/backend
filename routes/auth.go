package routes

import (
	"github.com/LitPad/backend/models"
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

	user := utils.ConvertStructData(data, models.User{}).(*models.User)
	// Validate email uniqueness
	db.Take(&user, models.User{Email: user.Email})
	if user.ID != uuid.Nil {
		data := map[string]string{
			"email": "Email already taken!",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Validate username uniqueness
	db.Take(&user, models.User{Username: user.Username})
	if user.ID != uuid.Nil {
		data := map[string]string{
			"username": "Username already taken!",
		}
		return c.Status(422).JSON(utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data))
	}

	// Create User
	db.Create(&user)

	// Send Email
	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	db.Save(&otp) // Create or save
	go senders.SendEmail(user, "activate", &otp.Code)

	response := schemas.RegisterResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Registration successful"}.Init(),
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

	user := models.User{Email: data.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_EMAIL, "Incorrect Email"))
	}

	if user.IsEmailVerified {
		return c.Status(200).JSON(schemas.ResponseSchema{Message: "Email already verified"}.Init())
	}

	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	if otp.ID == uuid.Nil || otp.Code != data.Otp {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_OTP, "Incorrect Otp"))
	}

	if otp.CheckExpiration() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_EXPIRED_OTP, "Expired Otp"))
	}

	// Update User
	user.IsEmailVerified = true
	db.Save(&user)

	// Send Welcome Email
	go senders.SendEmail(&user, "welcome", nil)
	response := schemas.ResponseSchema{Message: "Account verification successful"}.Init()
	return c.Status(200).JSON(response)
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
	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	db.Save(&otp) // Create or save
	go senders.SendEmail(&user, "activate", &otp.Code)

	response := schemas.ResponseSchema{Message: "Verification email sent"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Send Password Reset Otp
// @Description `This endpoint sends new password reset otp to the user's email.`
// @Tags Auth
// @Param email body schemas.EmailRequestSchema true "Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/send-password-reset-otp [post]
func (ep Endpoint) SendPasswordResetOtp(c *fiber.Ctx) error {
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
	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	db.Save(&otp) // Create or save
	go senders.SendEmail(&user, "reset", &otp.Code)

	response := schemas.ResponseSchema{Message: "Password otp sent"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Set New Password
// @Description `This endpoint verifies the password reset otp.`
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

	user := models.User{Email: data.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_EMAIL, "Incorrect Email"))
	}

	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	if otp.ID == uuid.Nil || otp.Code != data.Otp {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_INCORRECT_OTP, "Incorrect Otp"))
	}

	if otp.CheckExpiration() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_EXPIRED_OTP, "Expired Otp"))
	}

	// Set Password
	user.Password = utils.HashPassword(data.Password)
	db.Save(&user)

	// Send Email
	go senders.SendEmail(&user, "reset-success", nil)

	response := schemas.ResponseSchema{Message: "Password reset successful"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Login a user
// @Description This endpoint generates new access and refresh tokens for authentication
// @Tags Auth
// @Param user body schemas.LoginSchema true "User login"
// @Success 201 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Security GuestUserAuth
// @Router /auth/login [post]
func (ep Endpoint) Login(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.LoginSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := models.User{Email: data.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil || !utils.CheckPasswordHash(data.Password, user.Password) {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_CREDENTIALS, "Invalid Credentials"))
	}

	if !user.IsEmailVerified {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNVERIFIED_USER, "Verify your email first"))
	}

	// Create Auth Tokens
	access := GenerateAccessToken(user.ID)
	user.Access = &access
	refresh := GenerateRefreshToken()
	user.Refresh = &refresh
	db.Save(&user)
	response := schemas.LoginResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Login successful"}.Init(),
		Data:           schemas.TokensResponseSchema{Access: *user.Access, Refresh: *user.Refresh},
	}
	return c.Status(201).JSON(response)
}

// @Summary Refresh tokens
// @Description This endpoint refresh tokens by generating new access and refresh tokens for a user
// @Tags Auth
// @Param refresh body schemas.RefreshTokenSchema true "Refresh token"
// @Success 201 {object} schemas.ResponseSchema
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

	token := data.Refresh
	user := models.User{Refresh: &token}
	db.Take(&user, user)
	if user.ID == uuid.Nil || !DecodeRefreshToken(token){
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, "Refresh token is invalid or expired"))

	}

	// Create and Update Auth Tokens
	access := GenerateAccessToken(user.ID)
	user.Access = &access
	refresh := GenerateRefreshToken()
	user.Refresh = &refresh
	db.Save(&user)

	response := schemas.LoginResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Tokens refresh successful"}.Init(),
		Data:           schemas.TokensResponseSchema{Access: access, Refresh: refresh},
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
	db := ep.DB
	user := c.Locals("user").(*models.User)
	user.Access = nil
	user.Refresh = nil
	db.Save(user)
	response := schemas.ResponseSchema{Message: "Logout successful"}.Init()
	return c.Status(200).JSON(response)
}