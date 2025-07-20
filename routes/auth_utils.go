package routes

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	fb "github.com/huandu/facebook/v2"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AccessTokenPayload struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type RefreshTokenPayload struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user models.User, admin...bool) string {
	cfg := config.GetConfig()
	expiryMins := cfg.AccessTokenExpireMinutes
	if len(admin) > 0 {
		expiryMins = 35
	}
	expirationTime := time.Now().Add(time.Duration(expiryMins) * time.Minute)
	payload := AccessTokenPayload{
		UserId: user.ID,
		Username: user.Username,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(cfg.SecretKeyByte)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Access token: ", err)
	}
	return tokenString
}

func GenerateRefreshToken() string {
	cfg := config.GetConfig()

	expirationTime := time.Now().Add(time.Duration(cfg.RefreshTokenExpireMinutes) * time.Minute)
	payload := RefreshTokenPayload{
		Data: utils.GetRandomString(10),
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(cfg.SecretKeyByte)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Refresh token: ", err)
	}
	return tokenString
}

func DecodeAccessToken(token string, db *gorm.DB) (*models.User, *string) {
	cfg := config.GetConfig()

	claims := &AccessTokenPayload{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return cfg.SecretKeyByte, nil
	})
	tokenErr := "Auth Token is Invalid or Expired!"
	if err != nil {
		return nil, &tokenErr
	}
	if !tkn.Valid {
		return nil, &tokenErr
	}
	tokenObj := models.AuthToken{UserID: claims.UserId, Access: token}
	result := db.Joins("User").Take(&tokenObj, tokenObj)
	if result.Error != nil {
		return nil, &tokenErr
	}
	return &tokenObj.User, nil
}

func DecodeRefreshToken(token string) bool {
	cfg := config.GetConfig()

	claims := &RefreshTokenPayload{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return cfg.SecretKeyByte, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("JWT Error: ", "Invalid Signature")
		} else {
			log.Println("JWT Error: ", err)
		}
		return false
	}
	if !tkn.Valid {
		log.Println("Invalid Refresh Token")
		return false
	}
	return true
}

// Social Auth
// GOOGLE
type GooglePayload struct {
	SUB           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func ConvertGoogleToken(accessToken string, deviceType choices.DeviceType) (*GooglePayload, *utils.ErrorResponse) {
	cfg := config.GetConfig()

	clientID := cfg.GoogleAndroidClientID
	if deviceType == choices.DT_IOS {
		clientID = cfg.GoogleIOSClientID
	}
	payload, err := idtoken.Validate(context.Background(), accessToken, clientID)
	if err != nil {
		errMsg := "Invalid Token"
		if strings.Contains(err.Error(), "audience provided") {
			errMsg = "Invalid Audience"
		}
		errData := utils.RequestErr(utils.ERR_INVALID_TOKEN, errMsg)
		return nil, &errData
	}

	// Bind JSON into struct
	data := GooglePayload{}
	mapstructure.Decode(payload.Claims, &data)
	return &data, nil
}

// FACEBOOK
type FacebookPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ConvertFacebookToken(accessToken string) (*FacebookPayload, *utils.ErrorResponse) {
	cfg := config.GetConfig()

	res, err := fb.Get("/app", fb.Params{
		"access_token": accessToken,
	})
	if err != nil {
		errData := utils.RequestErr(utils.ERR_INVALID_TOKEN, "Token is invalid or expired")
		return nil, &errData
	}
	if res["id"] != cfg.FacebookAppID {
		errData := utils.RequestErr(utils.ERR_INVALID_TOKEN, "Invalid Facebook App ID")
		return nil, &errData
	}
	res, err = fb.Get("/me", fb.Params{
		"fields":       "id,name,email",
		"access_token": accessToken,
	})

	if err != nil {
		errData := utils.RequestErr(utils.ERR_INVALID_TOKEN, "Token is invalid or expired")
		return nil, &errData
	}
	// Bind JSON into struct
	data := FacebookPayload{}
	mapstructure.Decode(res, &data)
	return &data, nil
}

func RegisterSocialUser(db *gorm.DB, email string, name string, avatar *string, authType string) (*models.User, *models.AuthToken, *utils.ErrorResponse) {
	cfg := config.GetConfig()

	user := models.User{Email: email}
	db.Scopes(scopes.FollowerFollowingPreloaderScope).Take(&user, user)
	if user.ID == uuid.Nil {
		user = models.User{Name: &name, Email: email, IsEmailVerified: true, Password: cfg.SocialsPassword, Avatar: *avatar, SocialLogin: true}
		db.Create(&user)
	} else {
		if !user.SocialLogin {
			errData := utils.RequestErr(utils.ERR_INVALID_AUTH, fmt.Sprintf("This account wasn't created via %s. Please sign in using your email and password.", authType))
			return nil, nil, &errData
		}
	}
	// Generate tokens
	token := userManager.GenerateAuthTokens(db, user, GenerateAccessToken(user), GenerateRefreshToken())
	return &user, &token, nil
}
