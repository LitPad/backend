package managers

import (
	"fmt"
	"math"
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserManager struct {
	Model models.User
}

func (u UserManager) GetAll(db *gorm.DB, accountType *choices.AccType, staff *bool) []models.User {
	users := []models.User{}
	query := db
	if accountType != nil {
		query = query.Where(models.User{AccountType: *accountType})
	}
	if staff != nil {
		query = query.Where(models.User{IsStaff: *staff})
	}
	query.Scopes(scopes.FollowerFollowingBooksPreloaderScope).Order("created_at DESC").Find(&users)
	return users
}

func (u UserManager) GetCount(db *gorm.DB) int64 {
	var count int64
	db.Model(u.Model).Count(&count)
	return count
}

func (u UserManager) GetActiveSubscribersCount(db *gorm.DB) int64 {
	var count int64
	db.Model(u.Model).Where("subscription_expiry IS NOT NULL AND subscription_expiry > ?", time.Now()).Count(&count)
	return count
}

func (u UserManager) GetByUsername(db *gorm.DB, username string) *models.User {
	user := models.User{Username: username}
	db.Scopes(scopes.FollowerFollowingPreloaderScope).Take(&user, user)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (u UserManager) GetByEmail(db *gorm.DB, email string) *models.User {
	user := models.User{Email: email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (u UserManager) GetReaderByUsername(db *gorm.DB, username string) *models.User {
	user := models.User{Username: username, AccountType: choices.ACCTYPE_READER}
	db.Scopes(scopes.VerifiedUserScope).Take(&user, user)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (u UserManager) GetWriterByUsername(db *gorm.DB, username string) *models.User {
	user := models.User{Username: username, AccountType: choices.ACCTYPE_AUTHOR}
	db.Scopes(scopes.VerifiedUserScope).Take(&user, user)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (u UserManager) GetSubscribers(db *gorm.DB, subscriptionType *choices.SubscriptionTypeChoice, subscriptionStatus *choices.SubscriptionStatusChoice) []models.User {
	subscribers := []models.User{}
	query := db.Where("subscription_expiry IS NOT NULL")
	if subscriptionType != nil {
		query = query.Where("current_plan = ?", subscriptionType)
	}
	if subscriptionStatus != nil {
		currentTime := time.Now()
		if *subscriptionStatus == choices.SS_ACTIVE {
			query = query.Where("subscription_expiry > ?", currentTime)
		} else {
			query = query.Where("subscription_expiry < ?", currentTime)
		}
	}
	query.Find(&subscribers)
	return subscribers
}

func (u UserManager) GetUserPlanPercentages(db *gorm.DB) schemas.SubscriptionPlansAndPercentages {
	type Result struct {
		Category string
		Count    int64
	}

	var results []Result
	var totalUsers int64

	// Count the total users
	db.Model(&u.Model).Count(&totalUsers)

	// Count users by category (freeTier, monthly, annual)
	db.Model(&u.Model).
		Select(`CASE 
                    WHEN current_plan IS NULL THEN 'freeTier' 
                    WHEN current_plan = ? THEN 'monthly' 
                    WHEN current_plan = ? THEN 'annual' 
                END as category, COUNT(*) as count`,
			choices.ST_MONTHLY,
			choices.ST_ANNUAL).
		Group("category").
		Scan(&results)

	// Calculate percentages
	percentages := map[string]float64{
		"freeTier": 0,
		"monthly":  0,
		"annual":   0,
	}

	for _, result := range results {
		if totalUsers > 0 {
			percentages[result.Category] = math.Round((float64(result.Count) / float64(totalUsers)) * 100) // to 1 dp
		}
	}
	percentagesData := schemas.SubscriptionPlansAndPercentages{
		FreeTier: percentages["freeTier"],
		Monthly:  percentages["monthly"],
		Annual:   percentages["annual"],
	}
	return percentagesData
}

func (u UserManager) GetUserGrowthData(db *gorm.DB, choice choices.UserGrowthChoice) []schemas.UserGrowthData {
	var results []schemas.UserGrowthData
	var startDate time.Time
	var groupBy string

	// Calculate start date and grouping based on UserGrowthChoice
	now := time.Now()
	switch choice {
	case choices.UG_7:
		startDate = now.AddDate(0, 0, -7)
		groupBy = "DATE(created_at)" // Group by date
	case choices.UG_30:
		startDate = now.AddDate(0, 0, -30)
		groupBy = "DATE(created_at)" // Group by date
	case choices.UG_365:
		startDate = now.AddDate(0, 0, -365)
		groupBy = "DATE(created_at)" // Group by year-month
	}

	// Query the database
	db.Model(&u.Model).
		Select(fmt.Sprintf("%s AS period, COUNT(*) AS count", groupBy)).
		Where("created_at >= ?", startDate).
		Group(groupBy).
		Order("period ASC").
		Scan(&results)
	return results
}

type NotificationManager struct{}

func (n NotificationManager) GetAllByUser(db *gorm.DB, user *models.User) []models.Notification {
	notifications := []models.Notification{}
	db.Scopes(scopes.NotificationRelatedScope).Where(models.Notification{ReceiverID: user.ID}).Find(&notifications)
	return notifications
}

func (n NotificationManager) GetOneByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Notification {
	notification := models.Notification{ReceiverID: user.ID}
	db.Where("id = ?", id).Take(&notification, notification)
	if notification.ID == uuid.Nil {
		return nil
	}
	return &notification
}

func (n NotificationManager) MarkAsRead(db *gorm.DB, user *models.User) {
	db.Model(&models.Notification{}).Where("receiver_id = ?", user.ID).Updates(models.Notification{IsRead: true})
}

func (n NotificationManager) ReadOne(db *gorm.DB, user *models.User, id uuid.UUID) *utils.ErrorResponse {
	notification := models.Notification{ReceiverID: user.ID}
	db.Where("id = ?", id).Take(&notification, notification)
	if notification.ID == uuid.Nil {
		errD := utils.NotFoundErr("User has no notification with that ID")
		return &errD
	}
	notification.IsRead = true
	db.Save(&notification)
	return nil
}

func (n NotificationManager) Create(db *gorm.DB, sender *models.User, receiver models.User, ntype choices.NotificationTypeChoice, text string, book *models.Book, reviewID *uuid.UUID, replyID *uuid.UUID, sentGiftID *uuid.UUID) models.Notification {
	notification := models.Notification{
		SenderID: sender.ID, Sender: *sender, ReceiverID: receiver.ID,
		Ntype: ntype, ReviewID: reviewID, ReplyID: replyID, SentGiftID: sentGiftID, Text: text,
	}
	if book != nil {
		notification.Book = book
		notification.BookID = &book.ID
	}
	db.Create(&notification)
	return notification
}
