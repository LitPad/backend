package routes

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	fb "github.com/huandu/facebook/v2"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

var cfg = config.GetConfig()
var SECRETKEY = []byte(cfg.SecretKey)

type AccessTokenPayload struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshTokenPayload struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId uuid.UUID) string {
	expirationTime := time.Now().Add(time.Duration(cfg.AccessTokenExpireMinutes) * time.Minute)
	payload := AccessTokenPayload{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(SECRETKEY)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Access token: ", err)
	}
	return tokenString
}

func GenerateRefreshToken() string {
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
	tokenString, err := token.SignedString(SECRETKEY)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Refresh token: ", err)
	}
	return tokenString
}

func DecodeAccessToken(token string, db *gorm.DB) (*models.User, *string) {
	claims := &AccessTokenPayload{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRETKEY, nil
	})
	tokenErr := "Auth Token is Invalid or Expired!"
	if err != nil {
		return nil, &tokenErr
	}
	if !tkn.Valid {
		return nil, &tokenErr
	}
	user := models.User{BaseModel: models.BaseModel{ID: claims.UserId}, Access: &token}
	// Fetch Jwt model object
	result := db.Take(&user, user)
	if result.Error != nil {
		return nil, &tokenErr
	}
	return &user, nil
}

func DecodeRefreshToken(token string) bool {
	claims := &RefreshTokenPayload{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRETKEY, nil
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

func ConvertGoogleToken(accessToken string) (*GooglePayload, *utils.ErrorResponse) {
	payload, err := idtoken.Validate(context.Background(), accessToken, cfg.GoogleClientID)
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
		if err != nil {
			errData := utils.RequestErr(utils.ERR_INVALID_TOKEN, "Token is invalid or expired")
			return nil, &errData
		}
	}
	// Bind JSON into struct
	data := FacebookPayload{}
	mapstructure.Decode(res, &data)
	return &data, nil
}

func GenerateUsername(db *gorm.DB, firstName string, lastName string, username *string) string {
	uniqueUsername := slug.Make(firstName + " " + lastName)
	if username != nil {
		uniqueUsername = *username
	}
	user := models.User{Username: uniqueUsername}
	db.Take(&user, user)
	if user.ID != uuid.Nil {
		// username is already taken
		// Make it unique by attaching a random string
		// to it and repeat the function
		randomStr := utils.GetRandomString(6)
		uniqueUsername = uniqueUsername + "-" + randomStr
		return GenerateUsername(db, firstName, lastName, &uniqueUsername)
	}
	return uniqueUsername
}

func RegisterSocialUser(db *gorm.DB, email string, name string, avatar *string) (*models.User, *utils.ErrorResponse) {
	user := models.User{Email: email}
	db.Scopes(scopes.FollowerFollowingPreloaderScope).Take(&user, user)
	if user.ID == uuid.Nil {
		name := strings.Split(name, " ")
		firstName := name[0]
		lastName := name[1]
		username := GenerateUsername(db, firstName, lastName, nil)
		user = models.User{FirstName: firstName, LastName: lastName, Username: username, Email: email, IsEmailVerified: true, Password: utils.HashPassword(cfg.SocialsPassword), TermsAgreement: true, Avatar: avatar, SocialLogin: true}
		db.Create(&user)
	} else {
		if !user.SocialLogin {
			errData := utils.RequestErr(utils.ERR_INVALID_AUTH, "Requires password to login")
			return nil, &errData
		}
	}
	// Generate tokens
	access := GenerateAccessToken(user.ID)
	user.Access = &access
	refresh := GenerateRefreshToken()
	user.Refresh = &refresh
	db.Save(&user)
	return &user, nil
}
