package routes

import (
	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
)

var transactionManager = managers.TransactionManager{Model: models.Transaction{}}

// @Summary Admin Dashboard
// @Description `Retrieves minimal book data, counts and other metrics`
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_growth_filter query int64 false "User Growth to filter by in days" Enums(7, 30, 365)
// @Success 200 {object} schemas.DashboardResponseSchema "Successfully retrieved admin dashboard data"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin [get]
// @Security BearerAuth
func (ep Endpoint) AdminDashboard(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	userGrowthFilterQuery := c.QueryInt("user_growth_filter", int(choices.UG_365))
	userGrowthFilter := choices.UserGrowthChoice(userGrowthFilterQuery)
	if !userGrowthFilter.IsValid() {
		return c.Status(400).JSON(utils.InvalidParamErr("Invalid user growth filter choice!"))
	}
	totalUsers := userManager.GetCount(db)
	activeSubscribers := userManager.GetActiveSubscribersCount(db)
	subscriptionRevenue := transactionManager.GetSubscriptionRevenue(db)
	userSubscriptionPlanPercentages := userManager.GetUserPlanPercentages(db)
	userGrowthData := userManager.GetUserGrowthData(db, userGrowthFilter)
	books := bookManager.GetBooksOrderedByRatingAndVotes(db)
	response := schemas.DashboardResponseSchema{
		ResponseSchema: ResponseMessage("Admin Dashboard data retrieved successfully"),
		Data: schemas.DashboardResponseDataSchema{
			Username: user.Username, Avatar: user.Avatar,
			TotalUsers: totalUsers, ActiveSubscribers: activeSubscribers,
			SubscriptionRevenue: subscriptionRevenue, UserSubscriptionPlanPercentages: userSubscriptionPlanPercentages,
			UserGrowthData: userGrowthData, Books: books,
		},
	}
	return c.Status(200).JSON(response)
}
