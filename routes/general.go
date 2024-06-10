package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
)

// @Summary Retrieve site details
// @Description This endpoint retrieves few details of the site/application.
// @Tags General
// @Success 200 {object} schemas.SiteDetailResponseSchema
// @Router /general/site-detail [get]
func (ep Endpoint) GetSiteDetails(c *fiber.Ctx) error {
	db := ep.DB
	var sitedetail models.SiteDetail
	db.FirstOrCreate(&sitedetail, sitedetail)
	responseSiteDetail := schemas.SiteDetailResponseSchema{
		ResponseSchema: ResponseMessage("Site Details Fetched!"),
		Data:           sitedetail,
	}
	return c.Status(200).JSON(responseSiteDetail)
}

// @Summary Add a subscriber
// @Description This endpoint creates a newsletter subscriber in our application
// @Tags General
// @Param subscriber body models.Subscriber true "Subscriber object"
// @Success 201 {object} schemas.SubscriberResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /general/subscribe [post]
func (ep Endpoint) Subscribe(c *fiber.Ctx) error {
	db := ep.DB
	subscriber := models.Subscriber{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &subscriber); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Create subscriber
	db.FirstOrCreate(&subscriber, models.Subscriber{Email: subscriber.Email})

	responseSubscriber := schemas.SubscriberResponseSchema{
		ResponseSchema: ResponseMessage("Subscription successful!"),
		Data:           subscriber,
	}
	return c.Status(200).JSON(responseSubscriber)
}
