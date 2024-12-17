package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
)

// @Summary List Subscribers with Pagination
// @Description Retrieves a list of subscribers with support for pagination and optional filtering based on user subscription type or status.
// @Tags Admin | Subscribers
// @Accept json
// @Produce json
// @Param sub_type query string false "Subscription Type to filter by" Enums(MONTHLY, ANNUAL)
// @Param sub_status query string false "Subscription Status to filter by" Enums(ACTIVE, EXPIRED)
// @Param page query int false "Current page" default(1)
// @Success 200 {object} schemas.UserProfilesResponseSchema "Successfully retrieved list of user subs"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/subscribers [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetSubscribers(c *fiber.Ctx) error {
	db := ep.DB
	subType := c.Query("sub_type", "")
	var subscriptionType *choices.SubscriptionTypeChoice

	subStatus := c.Query("sub_status", "")
	var subscriptionStatus *choices.SubscriptionStatusChoice

	if subType == "" {
		subscriptionType = nil
	} else {
		subscriptionType = (*choices.SubscriptionTypeChoice)(&subType)
		if !subscriptionType.IsValid() {
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid subscription type"))
		}
	}

	if subStatus == "" {
		subscriptionStatus = nil
	} else {
		subscriptionStatus = (*choices.SubscriptionStatusChoice)(&subStatus)
		if !subscriptionStatus.IsValid() {
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid subscription status"))
		}
	}

	subscribers := userManager.GetSubscribers(db, subscriptionType, subscriptionStatus)
	// Paginate and return subscribers
	paginatedData, paginatedSubscribers, err := PaginateQueryset(subscribers, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	subscribers = paginatedSubscribers.([]models.User)
	response := schemas.UserProfilesResponseSchema{
		ResponseSchema: ResponseMessage("Subscribers fetched successfully"),
		Data: schemas.UserProfilesResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(subscribers),
	}
	return c.Status(200).JSON(response)
}