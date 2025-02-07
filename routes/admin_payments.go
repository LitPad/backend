package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Update A Plan Amount
// @Description This endpoint allows an admin to change the amount of a plan
// @Tags Wallet
// @Param plan body schemas.SubscriptionPlanSchema true "Plan data"
// @Success 200 {object} schemas.SubscriptionPlanResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /admin/payments/plans [put]
// @Security BearerAuth
func (ep Endpoint) UpdateSubscriptionPlan(c *fiber.Ctx) error {
	db := ep.DB
	data := schemas.SubscriptionPlanSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	plan := models.SubscriptionPlan{SubType: data.SubType}
	db.Take(&plan, plan)
	plan.Amount = data.Amount
	db.Save(&plan)
	response := schemas.SubscriptionPlanResponseSchema{
		ResponseSchema: ResponseMessage("Plan updated successfully"),
		Data:           schemas.SubscriptionPlanSchema{}.Init(plan),
	}
	return c.Status(200).JSON(response)
}

// @Summary Latest Transactions with Pagination
// @Description Retrieves a list of current transactions with support for pagination and optional filtering based on username.
// @Tags Admin | Payments
// @Accept json
// @Produce json
// @Param username query string false "Username to filter by"
// @Param page query int false "Current page" default(1)
// @Success 200 {object} schemas.TransactionsResponseSchema "Successfully retrieved list of transactions"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/payments/transactions [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetTransactions(c *fiber.Ctx) error {
	db := ep.DB
	username := GetQueryValue(c, "username") 
	transactions := []models.Transaction{}
	query := db
	if username != nil {
		user := models.User{Username: *username}
		db.Take(&user, user)
		if user.ID != uuid.Nil {
			query = query.Where("user_id = ?", user.ID)
		} else {
			// Ensure the query empties
			query = query.Where("user_id = ?", uuid.Nil)
		}
	}
	
	query.Joins("Coin").Joins("SubscriptionPlan").Order("created_at DESC").Find(&transactions)
	// Paginate and return transactions
	paginatedData, paginatedTransactions, err := PaginateQueryset(transactions, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	transactions = paginatedTransactions.([]models.Transaction)
	response := schemas.TransactionsResponseSchema{
		ResponseSchema: ResponseMessage("Transactions fetched successfully"),
		Data: schemas.TransactionsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(transactions),
	}
	return c.Status(200).JSON(response)
}