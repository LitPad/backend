package routes

import (
	"fmt"

	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
)

var (
	giftManager      = managers.GiftManager{}
	sendGiftManager      = managers.SentGiftManager{}
)

// @Summary View All Available Gifts
// @Description This endpoint shows a user gifts that can be sent
// @Tags Gifts
// @Success 200 {object} schemas.GiftsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /gifts [get]
func (ep Endpoint) GetAllGifts(c *fiber.Ctx) error {
	db := ep.DB
	gifts := giftManager.GetAll(db)

	response := schemas.GiftsResponseSchema{
		ResponseSchema: ResponseMessage("Gifts fetched successfully"),
	}.Init(gifts)
	return c.Status(200).JSON(response)
}

// @Summary Send Gift
// @Description This endpoint allows a user to send a gift
// @Tags Gifts
// @Param username path string true "Username of the writer"
// @Param gift_slug path string true "Slug of the gift being sent"
// @Success 201 {object} schemas.SentGiftResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /gifts/{username}/{gift_slug}/send/ [get]
// @Security BearerAuth
func (ep Endpoint) SendGift(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	writerUsername := c.Params("username")
	giftSlug := c.Params("gift_slug")

	writer := userManager.GetWriterByUsername(db, writerUsername)
	if writer == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No writer with that username"))
	}

	if user.ID == writer.ID {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_NOT_ALLOWED, "You can't send gifts to yourself"))
	}

	gift := giftManager.GetBySlug(db, giftSlug)
	if gift == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No available gift with that slug"))
	}

	if gift.Price > user.Coins {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INSUFFICIENT_COINS, "You have insufficient coins for that gift"))

	}

	// Send gift
	sentGift := sendGiftManager.Create(db, *gift, *user, *writer)
	
	// Create and send notification in socket
	notification := notificationManager.Create(
		db, user, *writer, choices.NT_GIFT, 
		fmt.Sprintf("%s sent you a gift.", user.FullName()),
		nil, nil, nil, &sentGift.ID,
	)
	SendNotificationInSocket(c, notification)

	response := schemas.SentGiftResponseSchema{
		ResponseSchema: ResponseMessage("Gift sent successfully"),
		Data:           schemas.SentGiftSchema{}.Init(sentGift),
	}
	return c.Status(201).JSON(response)
}

// @Summary View All Gifts Sent To A Writer
// @Description This endpoint allows a writer to view all gifts that was sent to him/her
// @Tags Gifts
// @Param page query int false "Current Page" default(1)
// @Param claimed query string false "Filter by claimed value: CLAIMED or NOT_CLAIMED "
// @Success 200 {object} schemas.SentGiftsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /gifts/sent [get]
// @Security BearerAuth
func (ep Endpoint) GetAllSentGifts(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	claimed := c.Query("claimed", "")
	var sentGifts []models.SentGift

	if claimed != "" {
		claimOpt := true
		if claimed == "NOT_CLAIMED" {
			claimOpt = false
		} else if claimed != "CLAIMED" {
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "Invalid claimed param"))
		}
		sentGifts = sendGiftManager.GetByWriter(db, *user, claimOpt)
	} else {
		sentGifts = sendGiftManager.GetByWriter(db, *user)
	}


	// Paginate and return sent gifts
	paginatedData, paginatedSentGifts, err := PaginateQueryset(sentGifts, c, 100)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	sentGifts = paginatedSentGifts.([]models.SentGift)
	response := schemas.SentGiftsResponseSchema{
		ResponseSchema: ResponseMessage("Gifts fetched successfully"),
		Data: schemas.SentGiftsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(sentGifts),
	}
	return c.Status(200).JSON(response)
}

// @Summary Claim Gift
// @Description This endpoint allows a writer to claim a gift
// @Tags Gifts
// @Param id path string true "ID of the sent gift (uuid)"
// @Success 200 {object} schemas.SentGiftResponseSchema	
// @Failure 400 {object} utils.ErrorResponse
// @Router /gifts/sent/{id}/claim [get]
// @Security BearerAuth
func (ep Endpoint) ClaimGift(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	sentGiftID := c.Params("id")
	parsedID := ParseUUID(sentGiftID)
	if parsedID == nil {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_PARAM, "You entered an invalid uuid"))
	}
	sentGift := sendGiftManager.GetByWriterAndID(db, *user, *parsedID)
	if sentGift == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "No gift with that ID was sent to you"))
	}

	if !sentGift.Claimed {
		// Claim gift
		user.Coins += sentGift.Gift.Price
		sentGift.Claimed = true
		db.Save(&user)
		db.Save(&sentGift)
	}

	response := schemas.SentGiftResponseSchema{
		ResponseSchema: ResponseMessage("Gift claimed successfully"),
		Data:           schemas.SentGiftSchema{}.Init(*sentGift),
	}
	return c.Status(200).JSON(response)
}