package schemas

import "github.com/LitPad/backend/models"

type UserProfilesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []UserProfile `json:"users"`
}

func (u UserProfilesResponseDataSchema) Init(users []models.User) UserProfilesResponseDataSchema {
	// Set Initial Data
	userItems := make([]UserProfile, 0)

	for _, user := range users {
		userItems = append(userItems, UserProfile{}.Init(user, nil))
	}
	u.Items = userItems
	return u
}

type UserProfilesResponseSchema struct {
	ResponseSchema
	Data UserProfilesResponseDataSchema `json:"data"`
}

type UserGrowthData struct {
	Period string `json:"period"` // e.g., "Jan 2025", "Week 1", etc.
	Count  int    `json:"count"`  // Number of new users
}

type SubscriptionPlansAndPercentages struct {
	FreeTier float64 `json:"free_tier"`
	Monthly  float64 `json:"monthly"`
	Annual   float64 `json:"annual"`
}

type DashboardResponseDataSchema struct {
	Username                        string                          `json:"username"`
	Avatar                          string                          `json:"avatar"`
	TotalUsers                      int64                           `json:"total_users"`
	ActiveSubscribers               int64                           `json:"active_subscribers"`
	SubscriptionRevenue             float64                         `json:"subscription_revenue"`
	UserSubscriptionPlanPercentages SubscriptionPlansAndPercentages `json:"user_subscription_plan_percentages"`
	UserGrowthData                  []UserGrowthData                `json:"user_growth_data"`
	Books                           []BookWithStats                 `json:"books"`
}

type DashboardResponseSchema struct {
	ResponseSchema
	Data DashboardResponseDataSchema `json:"data"`
}

type InviteAdminSchema struct {
	Email  string `json:"email" validate:"email,required"`
	Admin  bool   `json:"admin" validate:"required"`
	Author bool   `json:"author" validate:"required"`
}
