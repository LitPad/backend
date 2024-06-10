package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserManager struct{}

func (u UserManager) GetByUsername(db *gorm.DB, username string) *models.User {
	user := models.User{Username: username}
	db.Scopes(scopes.VerifiedUserScope).Take(&user, user)
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
	user := models.User{Username: username, AccountType: choices.ACCTYPE_WRITER}
	db.Scopes(scopes.VerifiedUserScope).Take(&user, user)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

type NotificationManager struct{}

func (n NotificationManager) GetAllByUser(db *gorm.DB, user *models.User) []models.Notification {
	notifications := []models.Notification{}
	db.Scopes(scopes.NotificationRelatedScope).Where(models.Notification{ReceiverID: user.ID}).Find(&notifications)
	return notifications
}

func (n NotificationManager) MarkAsRead(db *gorm.DB, user *models.User) {
	db.Model(&models.Notification{ReceiverID: user.ID}).Updates(models.Notification{IsRead: true})
}

func (n NotificationManager) ReadOne(db *gorm.DB, user *models.User, id uuid.UUID) *utils.ErrorResponse {
	notification := models.Notification{ReceiverID: user.ID}
	db.Where("id = ?", id).Take(&notification, notification)
	if notification.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no notification with that ID")
		return &errD
	}
	notification.IsRead = true
	db.Save(&notification)
	return nil
}
