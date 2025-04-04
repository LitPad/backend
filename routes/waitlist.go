package routes

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/senders"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Add to Waitlist
// @Description Adds a user to the waitlist.
// @Tags Waitlist
// @Accept json
// @Produce json
// @Param data body schemas.AddToWaitlist true "Waitlist data"
// @Success 200 {object} schemas.WaitlistResponseSchema "Successfully added to waitlist"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "Invalid Genre ID"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /waitlist [post]
// @Security BearerAuth
func (ep Endpoint) AddToWaitlist(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.AddToWaitlist{}

	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	genre := genreManager.GetBySlug(db, data.GenreSlug)

	if genre == nil {
		return c.Status(404).JSON(utils.NotFoundErr("Genre does not exist"))
	}

	waitlist := models.Waitlist{
		BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:      data.Name,
		Email:     data.Email,
		GenreID:   genre.ID,
	}

	db.Take(&waitlist, models.Waitlist{Email: waitlist.Email})

	var existingWaitlist models.Waitlist

	db.Where("email = ?", waitlist.Email).First(&existingWaitlist)
	if existingWaitlist.ID == uuid.Nil {
		db.Create(&waitlist)
	}

	go senders.AddEmailToBrevo(waitlist.Name, waitlist.Email)
	response := schemas.WaitlistResponseSchema{
		ResponseSchema: ResponseMessage("Added to waitlist successfully"),
	}
	return c.Status(200).JSON(response)
}
