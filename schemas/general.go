package schemas

import "github.com/LitPad/backend/models"

type SiteDetailResponseSchema struct {
	ResponseSchema
	Data			models.SiteDetail		`json:"data"`
}

type SubscriberResponseSchema struct {
	ResponseSchema
	Data			models.Subscriber		`json:"data"`
}