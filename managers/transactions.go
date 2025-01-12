package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"gorm.io/gorm"
)

type TransactionManager struct {
	Model    models.Transaction
}

func (t TransactionManager) GetSubscriptionRevenue (db *gorm.DB) float64 {
	var totalRevenue float64

	db.Model(&t.Model).
		Select("COALESCE(SUM(CAST(subscription_plans.amount AS NUMERIC)), 0)").
		Joins("JOIN subscription_plans ON subscription_plans.id = transactions.subscription_plan_id").
		Where("transactions.subscription_plan_id IS NOT NULL").
		Where("transactions.payment_status = ?", choices.PSSUCCEEDED).
		Scan(&totalRevenue)
	return totalRevenue
}

