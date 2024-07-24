package schemas

import "github.com/google/uuid"

type AddToWaitlist struct {
	Name string `json:"name"`
	Email string `json:"email"`
	GenreID uuid.UUID `json:"genre_id"`
}

type WaitlistResponseSchema struct {
	ResponseSchema
}