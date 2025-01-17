package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/scopes"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GiftManager struct {
	Model     models.Gift
	ModelList []models.Gift
}

func (g GiftManager) GetAll(db *gorm.DB) []models.Gift {
	gifts := g.ModelList
	db.Find(&gifts)
	return gifts
}

func (g GiftManager) GetBySlug(db *gorm.DB, slug string) *models.Gift {
	gift := models.Gift{Slug: slug}
	db.Take(&gift, gift)
	if gift.ID == uuid.Nil {
		return nil
	}
	return &gift
}

type SentGiftManager struct {
	Model     models.SentGift
	ModelList []models.SentGift
}

func (g SentGiftManager) GetByWriter(db *gorm.DB, writer models.User, claimed ...bool) []models.SentGift {
	sentGifts := g.ModelList
	query := db.Scopes(scopes.SentGiftRelatedScope).Where(models.SentGift{ReceiverID: writer.ID})
	if len(claimed) > 0 {
		query = query.Where(map[string]interface{}{"claimed": claimed[0]})
	} 
	query.Order("created_at DESC").Find(&sentGifts)
	return sentGifts
}

func (g SentGiftManager) GetByWriterAndID(db *gorm.DB, writer models.User, id uuid.UUID) *models.SentGift {
	sentGift := models.SentGift{ReceiverID: writer.ID}
	db.Scopes(scopes.SentGiftRelatedScope).Where("sent_gifts.id = ?", id).Take(&sentGift, sentGift)
	if sentGift.ID == uuid.Nil {
		return nil
	}
	return &sentGift
}

func (s SentGiftManager) Create(db *gorm.DB, gift models.Gift, sender models.User, receiver models.User) models.SentGift {
	sentGift := models.SentGift{
		SenderID: sender.ID, Sender: sender,
		ReceiverID: receiver.ID, Receiver: receiver,
		GiftID: gift.ID, Gift: gift,
	}
	db.Create(&sentGift)
	sender.Coins -= gift.Price
	sender.Lanterns += gift.Lanterns
	db.Save(&sender)
	return sentGift
}

func (s SentGiftManager) Process(db *gorm.DB, gift models.Gift, sender models.User, receiver models.User) models.SentGift{
	sentGift := models.SentGift {
		SenderID: sender.ID,
		Sender: sender,
		ReceiverID: receiver.ID,
		Receiver: receiver,
		GiftID: gift.ID, Gift: gift,
	}

	db.Create(&sentGift)
	return sentGift
}