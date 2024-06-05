package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
)

type GiftSchema struct {
	Name     string `json:"name" example:"Red rose"`
	Slug     string `json:"slug" example:"red-rose"`
	Price    int    `json:"price" example:"500"`
	Image    string `json:"image" example:"https://img.url"`
	Lanterns int    `json:"lanterns" example:"2"`
}

type GiftsResponseSchema struct {
	ResponseSchema
	Data []GiftSchema `json:"data"`
}

func (g GiftsResponseSchema) Init(gifts []models.Gift) GiftsResponseSchema {
	// Set Initial Data
	convertedGiftsData := utils.ConvertStructData(gifts, []GiftSchema{}).(*[]GiftSchema)
	g.Data = *convertedGiftsData
	return g
}

// -------------------------------------------------

type SentGiftSchema struct {
	ID 	uuid.UUID	`json:"id" example:"2b3bd817-135e-41bd-9781-33807c92ff40"`
	Sender     UserDataSchema `json:"sender"`
	Receiver     UserDataSchema `json:"receiver"`
	Item    GiftSchema    `json:"gift"`
	Claimed    bool `json:"claimed"`
	CreatedAt 	time.Time `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (s SentGiftSchema) Init(sentGift models.SentGift) SentGiftSchema {
	s.ID = sentGift.ID
	s.Sender = s.Sender.Init(sentGift.Sender)
	s.Receiver = s.Receiver.Init(sentGift.Receiver)
	convertItem := utils.ConvertStructData(sentGift.Gift, GiftSchema{}).(*GiftSchema)
	s.Item = *convertItem
	s.Claimed = sentGift.Claimed
	s.CreatedAt = sentGift.CreatedAt
	return s
}

type SentGiftsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []SentGiftSchema `json:"gifts"`
}

func (s SentGiftsResponseDataSchema) Init(sentGifts []models.SentGift) SentGiftsResponseDataSchema {
	// Set Initial Data
	sentGiftItems := s.Items
	for _, sentGift := range sentGifts {
		sentGiftItems = append(sentGiftItems, SentGiftSchema{}.Init(sentGift))
	}
	s.Items = sentGiftItems
	return s
}


type SentGiftsResponseSchema struct {
	ResponseSchema
	Data SentGiftsResponseDataSchema `json:"data"`
}

type SentGiftResponseSchema struct {
	ResponseSchema
	Data SentGiftSchema `json:"data"`
}