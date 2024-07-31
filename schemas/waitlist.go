package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/google/uuid"
)

type AddToWaitlist struct {
	Name string `json:"name"`
	Email string `json:"email"`
	GenreSlug string `json:"genre_slug"`
}

type WaitlistResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	GenreID   uuid.UUID `json:"genre_id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"created_at" example:"2024-06-05T02:32:34.462196+01:00"`
}

func (w WaitlistResponse) Init(waitlist models.Waitlist) WaitlistResponse{
	w = WaitlistResponse{
		ID: waitlist.ID,
		Name: waitlist.Name,
		Email: waitlist.Email,
		GenreID: waitlist.GenreID,
		GenreName: waitlist.Genre.Name,
		CreatedAt: waitlist.CreatedAt,
	}

	return w
}

func (wr WaitlistResponseDataSchema) Init(waitlists []models.Waitlist) WaitlistResponseDataSchema {
	list := wr.Items

	for _, waitlist := range waitlists {
		list = append(list, WaitlistResponse{}.Init(waitlist))
	}

	wr.Items =list
	return wr
}


type WaitlistResponseSchema struct {
	ResponseSchema
}

type WaitlistResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []WaitlistResponse `json:"waitlist"`
}

type WaitlistListResponseSchema struct {
	ResponseSchema
	Data WaitlistResponseDataSchema `json:"data"`
}