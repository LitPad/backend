package tests

import (
	"fmt"
	"log"
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
		log.Println("EXPIRY: ", user.OtpExpiry)

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

func TestAuth(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	baseUrl := "/api/v1/auth"

	// Run Auth Endpoint Tests
	register(t, app, baseUrl)
	resendVerificationEmail(t, app, db, baseUrl)
	verifyEmail(t, app, db, baseUrl)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
