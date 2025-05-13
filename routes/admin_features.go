package routes

import (
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
)

// @Summary List Featured Content
// @Description Retrieves a list of a featured content and optional filtering based on location.
// @Tags Admin | Featured
// @Accept json
// @Produce json
// @Param location query string false "Location"
// @Success 200 {object} schemas.FeaturedContentsResponseSchema "Successfully retrieved list of books"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/featured-contents [get]
// @Security BearerAuth
func (ep Endpoint) AdminGetFeaturedContents(c *fiber.Ctx) error {
	db := ep.DB
	locationQuery := c.Query("location", "")
	var location *choices.FeaturedContentLocationChoice
	if locationQuery == "" {
		location = nil
	} else {
		location_ := choices.FeaturedContentLocationChoice(locationQuery)
		location = &location_
	}
	if location != nil && !location.IsValid() {
		return c.Status(400).JSON(utils.InvalidParamErr("Invalid location parameter"))
	}
	contents := featuredContentManager.GetAll(db, location, nil)
	response := schemas.FeaturedContentsResponseSchema{
		ResponseSchema: ResponseMessage("Featured contents fetched successfully"),
	}.Init(contents)
	return c.Status(200).JSON(response)
}

// @Summary Add Featured Content
// @Description Add a featured content.
// @Tags Admin | Featured
// @Accept json
// @Produce json
// @Param data body schemas.FeaturedContentEntrySchema true "content"
// @Success 200 {object} schemas.FeaturedContentResponseSchema "Featured content added successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/featured-contents [post]
// @Security BearerAuth
func (ep Endpoint) AdminAddAFeaturedContent(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.FeaturedContentEntrySchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	book, err := bookManager.GetBySlug(db, data.BookSlug, false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	content := featuredContentManager.Create(db, data.Location, data.Desc, *book)
	response := schemas.FeaturedContentResponseSchema{
		ResponseSchema: ResponseMessage("Featured Content Added Successfully"),
		Data: schemas.FeaturedContentSchema{}.Init(content),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update Featured Content
// @Description Update a featured content.
// @Tags Admin | Featured
// @Accept json
// @Produce json
// @Param id path string true "Featured Content ID (uuid)"
// @Param data body schemas.FeaturedContentEntrySchema true "content"
// @Success 200 {object} schemas.FeaturedContentResponseSchema "Featured content updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/featured-contents/{id} [put]
// @Security BearerAuth
func (ep Endpoint) AdminUpdateAFeaturedContent(c *fiber.Ctx) error {
	db := ep.DB
	id := ParseUUID(c.Params("id"))
	if id == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("Enter a valid uuid"))
	}
	content := featuredContentManager.GetByID(db, *id)
	if content == nil {
		return c.Status(404).JSON(utils.NotFoundErr("No featured content with that ID"))
	}
	data := schemas.FeaturedContentEntrySchema{}
	errCode, errData := ValidateRequest(c, &data)
	if errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	book, err := bookManager.GetBySlug(db, data.BookSlug, false)
	if err != nil {
		return c.Status(404).JSON(err)
	}
	updatedContent := featuredContentManager.Update(db, *content, data.Location, data.Desc, *book)
	response := schemas.FeaturedContentResponseSchema{
		ResponseSchema: ResponseMessage("Featured Content Updated Successfully"),
		Data: schemas.FeaturedContentSchema{}.Init(updatedContent),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete Featured Content
// @Description Delete a featured content.
// @Tags Admin | Featured
// @Accept json
// @Produce json
// @Param id path string true "Featured Content ID (uuid)"
// @Success 200 {object} schemas.ResponseSchema "Featured content deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/featured-contents/{id} [delete]
func (ep Endpoint) AdminDeleteAFeaturedContent(c *fiber.Ctx) error {
	db := ep.DB
	id := ParseUUID(c.Params("id"))
	if id == nil {
		return c.Status(400).JSON(utils.InvalidParamErr("Enter a valid uuid"))
	}
	content := featuredContentManager.GetByID(db, *id)
	if content == nil {
		return c.Status(404).JSON(utils.NotFoundErr("No featured content with that ID"))
	}
	db.Delete(&content)
	return c.Status(200).JSON(ResponseMessage("Featured Content Deleted Successfully"))
}
