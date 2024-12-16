package routes

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/senders"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

func generateTokenForWalletReq(conf config.Config) (string, error){
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": 80 * time.Second,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(conf.WalletSecret))
}

// const WALLET_SERVER = "http://backend:2500/api"
const WALLET_SERVER = "http://localhost:2500/api"


// @Summary Create a new ICP wallet
// @Description `This endpoint creates a new ICP wallet`
// @Tags Wallet
// @Param user body schemas.CreateICPWallet true "User data"
// @Failure 400 {object} utils.ErrorResponse
// @Router /wallet/icp [post]
func (ep Endpoint) CreateICPWallet(c *fiber.Ctx) error {
	endpoint := WALLET_SERVER + "/wallet"
	
	conf := config.GetConfig()

	token, err := generateTokenForWalletReq(conf)

	if err != nil{
		return c.Status(500).JSON(fiber.Map{"err": err})
	}

	headers := map[string]string{
		"Accept": "application/json",
		"Access": "Litpad " + token,
	}

	data := schemas.CreateICPWallet{}

	if errCode, errData := ValidateRequest(c, &data); errData != nil{
		return c.Status(*errCode).JSON(errData)
	}
	
	requestBody := []byte(fmt.Sprintf(`{"username": "%s"}`, data.Username))

	response, err := senders.MakeRequest(fasthttp.MethodPost, endpoint, headers, requestBody)

	if err != nil {
		return c.Status(response.StatusCode()).JSON(fiber.Map{"error": err.Error()})
	}

	var decodedResponse map[string]interface{}
	if err := json.Unmarshal(response.Body(), &decodedResponse); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse response"})
	}

	return c.Status(201).JSON(decodedResponse)
}


// @Summary Get user ICP wallet balance
// @Description This endpoint returns user ICP wallet balance
// @Tags Wallet
// @Param username path string true "Username of user"
// @Failure 400 {object} utils.ErrorResponse
// @Router /wallet/icp/{username}/balance [get]
func(ep Endpoint) GetICPWalletBalance(c *fiber.Ctx) error {
	username := c.Params("username")

	if username == "" {
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path params"))
	}

	endpoint := WALLET_SERVER + "/wallet/" + username + "/balance";

	conf := config.GetConfig()

	token, err := generateTokenForWalletReq(conf)

	if err != nil{
		return c.Status(500).JSON(fiber.Map{"err": err})
	}

	headers := map[string]string{
		"Accept": "application/json",
		"Access": "Litpad " + token,
	}

	response, err := senders.MakeRequest(fasthttp.MethodGet, endpoint, headers, nil)

	if err != nil{
				return c.Status(response.StatusCode()).JSON(fiber.Map{"error": err.Error()})
	}

	var decodedResponse map[string]interface{}

	if err := json.Unmarshal(response.Body(), &decodedResponse); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse response"})
	}

	return c.Status(201).JSON(decodedResponse)
}