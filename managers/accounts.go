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

func (u UserManager) GetByUsername(db *gorm.DB, username string) *models.User {
	user := models.User{Username: username}
	db.Scopes(scopes.FollowerFollowingPreloaderScope).Take(&user, user)
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
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "User has no notification with that ID")
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
