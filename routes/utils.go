package routes

import (
	"fmt"
	"mime/multipart"
	"strings"
	"net/http"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"gorm.io/gorm"
)

func ResponseMessage(message string) schemas.ResponseSchema {
	return schemas.ResponseSchema{Status: "success", Message: message}
}

func RequestUser(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func GetBaseReferer(c *fiber.Ctx) string {
	referer := c.Context().Referer()
	return string(referer[:])
}

func CreateCheckoutSession(c *fiber.Ctx, db *gorm.DB, user models.User, coin models.Coin, quantity int64) (*models.Transaction, *utils.ErrorResponse) {
	baseUrl := GetBaseReferer(c)
	stripe.Key = cfg.StripeSecretKey
	price := coin.Price.Mul(decimal.NewFromFloat(100)).IntPart()
	productName := fmt.Sprintf("%s coins", fmt.Sprint(coin.Amount))
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(baseUrl + cfg.StripeCheckoutSuccessUrlPath),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: &productName,
					},
					TaxBehavior: stripe.String(string(stripe.PriceTaxBehaviorExclusive)),
					UnitAmount:  stripe.Int64(price),
				},
				Quantity: stripe.Int64(quantity),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		CustomerEmail: &user.Email,
	}
	s, err := session.New(params)
	if err != nil {
		errD := utils.RequestErr(utils.ERR_SERVER_ERROR, "Something went wrong")
		return nil, &errD
	}

	// Create Transaction Object
	transaction := models.Transaction{Reference: s.ID, UserID: user.ID, CoinID: coin.ID, PaymentType: choices.PTYPE_STRIPE, Quantity: quantity, CheckoutURL: s.URL}
	db.Create(&transaction)
	transaction.Coin = coin
	return &transaction, nil
}

func IsValidPaymentStatus(s string) bool {
	switch choices.PaymentStatus(s) {
	case choices.PSPENDING, choices.PSSUCCEEDED, choices.PSFAILED, choices.PSCANCELED:
		return true
	}
	return false
}

func ValidatePaymentStatus(c *fiber.Ctx) (*string, *utils.ErrorResponse) {
	status := c.Query("payment_status", "")
	if status != "" && !IsValidPaymentStatus(status) {
		errD := utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid payment status")
		return nil, &errD
	}
	return &status, nil
}

func CheckTagStrings(db *gorm.DB, submittedList []string) ([]models.Tag, *string) {
	tags := []models.Tag{}
	db.Find(&tags)
    // Create a map for quick lookup of predefined strings
    predefinedMap := make(map[string]bool)
    for _, item := range tags {
        predefinedMap[item.Slug] = true
    }

    // Iterate over the submitted list and check for any missing strings
    missingStrings := []string{}
    for _, item := range submittedList {
        if !predefinedMap[item] {
            missingStrings = append(missingStrings, item)
        }
    }

    // Return a message based on the result
    if len(missingStrings) > 0 {
		missingTags := strings.Join(missingStrings, ", ")
        errMsg := fmt.Sprintf("The following are invalid tag slugs: %v", missingTags)
		return tags, &errMsg
    }
	tagsToReturn := []models.Tag{}
	db.Where("slug IN ?", submittedList).Find(&tagsToReturn)
	return tagsToReturn, nil
}

func ValidateAndUploadImage(c *fiber.Ctx, name string, folder string, required bool) (*string, *utils.ErrorResponse) {
	file, err := c.FormFile(name)

	data := map[string]string{
		name: "Invalid image type",
	}
	errData := utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)

	if required && err != nil {
		data[name] = "Image is required"
		errData = utils.RequestErr(utils.ERR_INVALID_ENTRY, "Invalid Entry", data)
		return nil, &errData
	}

	// Open the file
	if file != nil {
		fileHandle, err := file.Open()
		if err != nil {
			return nil, &errData
		}
		
		defer fileHandle.Close()

		// Read the first 512 bytes for content type detection
		buffer := make([]byte, 512)
		_, err = fileHandle.Read(buffer)
		if err != nil {
			return nil, &errData
		}

		// Detect the content type
		contentType := http.DetectContentType(buffer)
		switch contentType {
			case "image/jpeg", "image/png", "image/gif":
				// Upload file
				fileUrl := uploadToCloudinary(c, file, folder)
				return &fileUrl, nil
		}
		return nil, &errData
	}
	return nil, nil
}

// uploadToCloudinary uploads the file to Cloudinary and returns the URL of the uploaded file
func uploadToCloudinary(c *fiber.Ctx, file *multipart.FileHeader, folder string) string {
	directory := fmt.Sprintf("litpad/%s", folder)
	// Set up Cloudinary
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryApiKey, cfg.CloudinaryApiSecret)
	if err != nil {
		return ""
	}

	// Open the file for reading
	fileHandle, err := file.Open()
	if err != nil {
		return ""
	}
	defer fileHandle.Close()

	// Upload the file to Cloudinary
	uploadParams := uploader.UploadParams{Folder: directory}
	uploadResult, err := cld.Upload.Upload(c.Context(), fileHandle, uploadParams)
	if err != nil {
		return ""
	}

	return uploadResult.SecureURL
}