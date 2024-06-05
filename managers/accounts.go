package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type UserManager struct {}

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