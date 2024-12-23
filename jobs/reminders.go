package jobs

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/senders"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type EmailTaskPayload struct {
	UserID      uuid.UUID
	EmailType   senders.EmailTypeChoice
	TokenString *string
	Url         *string
	ExtraData   map[string]interface{}
}

const TypeSendEmail = "send_email"

// EmailTaskHandler handles tasks for sending emails.
func EmailTaskHandler(db *gorm.DB) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload EmailTaskPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			log.Printf("Error unmarshaling task payload: %v\n", err)
			return err
		}

		// Retrieve user from database
		var user models.User
		if err := db.First(&user, payload.UserID).Error; err != nil {
			log.Printf("Error finding user with ID %d: %v\n", payload.UserID, err)
			return err
		}

		// Call the SendEmail function
		senders.SendEmail(&user, payload.EmailType, nil, nil, payload.ExtraData)
		return nil
	}
}

func QueueEmailTask(redisClient *asynq.Client, user *models.User, emailType senders.EmailTypeChoice, tokenString *string, url *string, extraData map[string]interface{}) {
	payload := EmailTaskPayload{
		UserID:      user.ID,
		EmailType:   emailType,
		TokenString: tokenString,
		Url:         url,
		ExtraData:   extraData,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling email task payload: %v\n", err)
		return
	}

	task := asynq.NewTask(TypeSendEmail, data)
	if _, err := redisClient.Enqueue(task, asynq.Queue("critical")); err != nil {
		log.Printf("Failed to enqueue email task: %v\n", err)
	}
}

func ReminderJob(db *gorm.DB, redisClient *asynq.Client) {
	currentTime := time.Now()
	oneWeekLater := currentTime.Add(7 * 24 * time.Hour)

	// Find users with subscription expiring soon
	var expiringUsers []models.User
	db.Where("subscription_expiry BETWEEN ? AND ? AND reminder_sent = ?", currentTime, oneWeekLater, false).Find(&expiringUsers)

	// Send reminders for expiring subscriptions via Asynq tasks
	for _, user := range expiringUsers {
		// Queue a task for sending subscription-expiring email
		extraData := map[string]interface{}{"subscriptionType": user.CurrentPlan}
		QueueEmailTask(redisClient, &user, senders.ET_SUBSCRIPTION_EXPIRING, nil, nil, extraData)
	}

	// Bulk update reminder_sent for users with expiring subscriptions
	if len(expiringUsers) > 0 {
		var userIds []uuid.UUID
		for _, user := range expiringUsers {
			userIds = append(userIds, user.ID)
		}
		db.Model(&models.User{}).Where("id IN ?", userIds).Update("reminder_sent", true)
	}

	// Find users with expired subscriptions
	var expiredUsers []models.User
	db.Where("subscription_expiry < ? AND reminder_sent = ?", currentTime, false).Find(&expiredUsers)

	// Send reminders for expired subscriptions via Asynq tasks
	for _, user := range expiredUsers {
		// Queue a task for sending subscription-expired email
		extraData := map[string]interface{}{"subscriptionType": user.CurrentPlan}
		QueueEmailTask(redisClient, &user, senders.ET_SUBSCRIPTION_EXPIRED, nil, nil, extraData)
	}

	// Bulk update reminder_sent for users with expired subscriptions
	if len(expiredUsers) > 0 {
		var userIds []uuid.UUID
		for _, user := range expiredUsers {
			userIds = append(userIds, user.ID)
		}
		db.Model(&models.User{}).Where("id IN ?", userIds).Update("reminder_sent", true)
	}
}