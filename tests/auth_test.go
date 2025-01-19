package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/LitPad/backend/database"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func register(t *testing.T, app *fiber.App, baseUrl string) {
	t.Run("Register User Successfully", func(t *testing.T) {
		url := fmt.Sprintf("%s/register", baseUrl)
		validEmail := "testregisteruser@email.com"
		userData := schemas.RegisterUser{
			Email:    validEmail,
			Password: "testregisteruserpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", userData)

		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Registration successful", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})

	t.Run("Register User Failure By Already Used details", func(t *testing.T) {
		email := "testregisteruser@email.com"
		url := fmt.Sprintf("%s/register", baseUrl)
		userData := schemas.RegisterUser{
			Email:    email,
			Password: "testregisteruserpassword",
		}

		// Verify that a user with the same email cannot be registered again
		res := ProcessTestBody(t, app, url, "POST", userData)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_ENTRY, body["code"])
		assert.Equal(t, "Invalid Entry", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = "Email already taken!"
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})
}

func verifyEmail(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject verification due to invalid email or otp", func(t *testing.T) {
		user := TestUser(db)
		url := fmt.Sprintf("%s/verify-email", baseUrl)
		verificationData := schemas.VerifyEmailRequestSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: user.Email},
			Otp:                111111,
		}
		res := ProcessTestBody(t, app, url, "POST", verificationData)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Email or OTP", body["message"])
	})

	t.Run("Reject verification due to expired otp", func(t *testing.T) {
		user := TestUser(db)
		user.GenerateOTP(db)
		otpExpiry := time.Now().UTC().Add(-900 * time.Second)
		user.OtpExpiry = &otpExpiry
		db.Save(&user)

		url := fmt.Sprintf("%s/verify-email", baseUrl)
		verificationData := schemas.VerifyEmailRequestSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: user.Email},
			Otp:                *user.Otp,
		}
		res := ProcessTestBody(t, app, url, "POST", verificationData)

		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Expired OTP", body["message"])
	})

	t.Run("Accept verification due to correct entry", func(t *testing.T) {
		user := TestUser(db)
		user.GenerateOTP(db)
		db.Save(&user)

		url := fmt.Sprintf("%s/verify-email", baseUrl)
		verificationData := schemas.VerifyEmailRequestSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: user.Email},
			Otp:                *user.Otp,
		}
		res := ProcessTestBody(t, app, url, "POST", verificationData)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Account verification successful", body["message"])
	})
}

func resendVerificationEmail(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject verification due to invalid email", func(t *testing.T) {
		url := fmt.Sprintf("%s/resend-verification-email", baseUrl)
		emailRequestData := schemas.EmailRequestSchema{
			Email: "invalid@example.com",
		}
		res := ProcessTestBody(t, app, url, "POST", emailRequestData)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Email", body["message"])
	})

	t.Run("Reject email resend due to already verified user but with success response", func(t *testing.T) {
		user := TestVerifiedUser(db)
		url := fmt.Sprintf("%s/resend-verification-email", baseUrl)
		emailRequestData := schemas.EmailRequestSchema{Email: user.Email}
		res := ProcessTestBody(t, app, url, "POST", emailRequestData)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Email already verified", body["message"])
	})

	t.Run("Accept email resend due to valid email", func(t *testing.T) {
		user := TestUser(db)
		url := fmt.Sprintf("%s/resend-verification-email", baseUrl)
		emailRequestData := schemas.EmailRequestSchema{Email: user.Email}
		res := ProcessTestBody(t, app, url, "POST", emailRequestData)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Verification email sent", body["message"])
	})
}

func sendPasswordResetLink(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject email sending due to invalid email", func(t *testing.T) {
		url := fmt.Sprintf("%s/send-password-reset-link", baseUrl)
		emailRequestData := schemas.EmailRequestSchema{
			Email: "invalid@example.com",
		}
		res := ProcessTestBody(t, app, url, "POST", emailRequestData)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Email", body["message"])
	})

	t.Run("Accept link sending due to valid email", func(t *testing.T) {
		user := TestUser(db)
		url := fmt.Sprintf("%s/send-password-reset-link", baseUrl)
		emailRequestData := schemas.EmailRequestSchema{Email: user.Email}
		res := ProcessTestBody(t, app, url, "POST", emailRequestData)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password reset link sent", body["message"])
	})
}

func verifyPasswordResetToken(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject verification due to invalid token", func(t *testing.T) {
		url := fmt.Sprintf("%s/verify-password-reset-token/invalid-token-string", baseUrl)
		res := ProcessTestGet(app, url)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Token", body["message"])
	})

	t.Run("Reject verification due to expired token", func(t *testing.T) {
		user := TestVerifiedUser(db)
		user.GenerateToken(db)
		tokenExpiry := time.Now().UTC().Add(-900 * time.Second)
		user.TokenExpiry = &tokenExpiry
		db.Save(&user)

		url := fmt.Sprintf("%s/verify-password-reset-token/%s", baseUrl, *user.TokenString)
		res := ProcessTestGet(app, url)

		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Expired Token", body["message"])
	})

	t.Run("Accept verification due to valid token", func(t *testing.T) {
		user := TestUser(db)
		user.GenerateToken(db)
		db.Save(&user)

		url := fmt.Sprintf("%s/verify-password-reset-token/%s", baseUrl, *user.TokenString)
		res := ProcessTestGet(app, url)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Token verified successfully", body["message"])
	})
}

func setNewPassword(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject password reset due to invalid email or token", func(t *testing.T) {
		url := fmt.Sprintf("%s/set-new-password", baseUrl)
		passwordResetData := schemas.SetNewPasswordSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: "invalid@example.com"},
			TokenString: "invalidtoken", Password: "newpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Email or Token", body["message"])
	})

	t.Run("Reject password reset due to expired token", func(t *testing.T) {
		user := TestVerifiedUser(db)
		user.GenerateToken(db)
		tokenExpiry := time.Now().UTC().Add(-900 * time.Second)
		user.TokenExpiry = &tokenExpiry
		db.Save(&user)

		url := fmt.Sprintf("%s/set-new-password", baseUrl)
		passwordResetData := schemas.SetNewPasswordSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: user.Email},
			TokenString: *user.TokenString, Password: "newpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Expired Token", body["message"])
	})

	t.Run("Accept password reset due to valid email and token", func(t *testing.T) {
		user := TestUser(db)
		user.GenerateToken(db)
		db.Save(&user)

		url := fmt.Sprintf("%s/set-new-password", baseUrl)
		passwordResetData := schemas.SetNewPasswordSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: user.Email},
			TokenString: *user.TokenString, Password: "newpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password reset successful", body["message"])
	})
}

func login(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject login due to invalid credentials", func(t *testing.T) {
		url := fmt.Sprintf("%s/login", baseUrl)
		loginData := schemas.LoginSchema{
			Email: "invalid@example.com",
			Password: "invalidpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", loginData)

		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Credentials", body["message"])
	})

	t.Run("Reject login due to unverified email", func(t *testing.T) {
		user := TestUser(db)
		url := fmt.Sprintf("%s/login", baseUrl)
		loginData := schemas.LoginSchema{
			Email: user.Email,
			Password: "testpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", loginData)

		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Verify your email first", body["message"])
	})

	t.Run("Accept login due to valid credentials", func(t *testing.T) {
		user := TestVerifiedUser(db)

		url := fmt.Sprintf("%s/login", baseUrl)
		loginData := schemas.LoginSchema{
			Email: user.Email,
			Password: "testpassword",
		}
		res := ProcessTestBody(t, app, url, "POST", loginData)

		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Login successful", body["message"])
	})
}

func googleLogin(t *testing.T, app *fiber.App, baseUrl string) {
	t.Run("Reject google login due to invalid token", func(t *testing.T) {
		url := fmt.Sprintf("%s/google", baseUrl)
		googleLoginData := schemas.SocialLoginSchema{Token: "invalid_token"}
		res := ProcessTestBody(t, app, url, "POST", googleLoginData)

		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Token", body["message"])
	})
}

func facebookLogin(t *testing.T, app *fiber.App, baseUrl string) {
	t.Run("Reject facebook login due to invalid token", func(t *testing.T) {
		url := fmt.Sprintf("%s/facebook", baseUrl)
		facebookLoginData := schemas.SocialLoginSchema{Token: "invalid_token"}
		res := ProcessTestBody(t, app, url, "POST", facebookLoginData)

		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Token is invalid or expired", body["message"])
	})
}

func refresh(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject Token Refresh Due To Invalid/Expired Token", func(t *testing.T) {
		url := fmt.Sprintf("%s/refresh", baseUrl)
		refreshData := schemas.RefreshTokenSchema{Refresh: "invalid"}
		res := ProcessTestBody(t, app, url, "POST", refreshData)

		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Refresh token is invalid or expired", body["message"])
	})

	t.Run("Accept Token Refresh Due To Valid Refresh Token", func(t *testing.T) {
		user := TestVerifiedUser(db)
		token := JwtData(db, user)

		url := fmt.Sprintf("%s/refresh", baseUrl)
		refreshData := schemas.RefreshTokenSchema{Refresh: *token.Refresh}
		res := ProcessTestBody(t, app, url, "POST", refreshData)

		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tokens refresh successful", body["message"])
	})
}

func logout(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject Logout Due To Invalid Token", func(t *testing.T) {
		url := fmt.Sprintf("%s/logout", baseUrl)
		res := ProcessTestGet(app, url, "invalid_token")
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Auth Token is Invalid or Expired!", body["message"])
	})

	t.Run("Accept Logout Due To Valid Token", func(t *testing.T) {
		url := fmt.Sprintf("%s/logout", baseUrl)
		token := AccessToken(db, TestVerifiedUser(db))
		res := ProcessTestGet(app, url, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Logout successful", body["message"])
	})
}

func TestAuth(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	baseUrl := "/api/v1/auth"

	// Run Auth Endpoint Tests
	register(t, app, baseUrl)
	resendVerificationEmail(t, app, db, baseUrl)
	verifyEmail(t, app, db, baseUrl)
	sendPasswordResetLink(t, app, db, baseUrl)
	verifyPasswordResetToken(t, app, db, baseUrl)
	setNewPassword(t, app, db, baseUrl)
	login(t, app, db, baseUrl)
	googleLogin(t, app, baseUrl)
	facebookLogin(t, app, baseUrl)
	refresh(t, app, db, baseUrl)
	logout(t, app, db, baseUrl)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
