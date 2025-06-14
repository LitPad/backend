package tests

import (
	"fmt"
	"testing"

	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func addBookGenre(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	genre := GenreData(db)
	genreData := schemas.TagsAddSchema{
		Name: genre.Name,
	}
	t.Run("Reject Genre Creation Due To Already Existing Genre", func(t *testing.T) {
		url := fmt.Sprintf("%s/genres", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", genreData, token)
		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Genre already exists", body["data"].(map[string]interface{})["name"])
	})

	t.Run("Reject Genre Creation Due To Invalid Tag Slug", func(t *testing.T) {
		genreData.Name = "Different Genre"
		url := fmt.Sprintf("%s/genres", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", genreData, token)
		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Enter at least one valid tag slug", body["data"].(map[string]interface{})["tag_slugs"])
	})

	t.Run("Accept Genre Creation", func(t *testing.T) {
		url := fmt.Sprintf("%s/genres", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", genreData, token)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Genre added successfully", body["message"])
	})
}

func addBookTag(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	tag := TagData(db)
	tagData := schemas.TagsAddSchema{Name: tag.Name}
	t.Run("Reject Tag Creation Due To Already Existing Tag", func(t *testing.T) {
		url := fmt.Sprintf("%s/tags", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", tagData, token)
		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Tag already exists", body["data"].(map[string]interface{})["name"])
	})

	t.Run("Accept Tag Creation", func(t *testing.T) {
		tagData.Name = "Different tag"
		url := fmt.Sprintf("%s/tags", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", tagData, token)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tag added successfully", body["message"])
	})
}

func updateBookGenre(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	genre := GenreData(db)
	genreData := schemas.TagsAddSchema{
		Name: genre.Name,
	}

	t.Run("Accept Genre Update", func(t *testing.T) {
		url := fmt.Sprintf("%s/genres/%s", baseUrl, genre.Slug)
		res := ProcessJsonTestBody(t, app, url, "PUT", genreData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Genre updated successfully", body["message"])
	})
}

func updateBookTag(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	tag := TagData(db)
	tagData := schemas.TagsAddSchema{Name: tag.Name}

	t.Run("Accept Tag Update", func(t *testing.T) {
		url := fmt.Sprintf("%s/tags/%s", baseUrl, tag.Slug)
		res := ProcessJsonTestBody(t, app, url, "PUT", tagData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tag updated successfully", body["message"])
	})
}

func deleteBookGenre(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	genre := GenreData(db)

	t.Run("Accept Genre Delete", func(t *testing.T) {
		url := fmt.Sprintf("%s/genres/%s", baseUrl, genre.Slug)
		res := ProcessTestGetOrDelete(app, url, "DELETE", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Genre deleted successfully", body["message"])
	})
}

func deleteBookTag(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	tag := TagData(db)

	t.Run("Accept Tag Delete", func(t *testing.T) {
		url := fmt.Sprintf("%s/tags/%s", baseUrl, tag.Slug)
		res := ProcessTestGetOrDelete(app, url, "DELETE", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tag deleted successfully", body["message"])
	})
}

func adminGetBooks(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	author := TestAuthor(db)
	BookData(db, author)

	t.Run("Accept Books Fetch", func(t *testing.T) {
		res := ProcessTestGetOrDelete(app, baseUrl, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Books fetched successfully", body["message"])
	})
}

func adminGetBooksByAuthor(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	author := TestAuthor(db)
	BookData(db, author)

	t.Run("Reject Books Fetch By Author Due To Invalid Username", func(t *testing.T) {
		url := fmt.Sprintf("%s/by-username/invalid-username", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Author does not exist!", body["message"])
	})

	t.Run("Accept Books Fetch By Author With Valid Username", func(t *testing.T) {
		url := fmt.Sprintf("%s/by-username/%s", baseUrl, author.Username)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Books fetched successfully", body["message"])
	})
}

func adminGetBookDetails(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	author := TestAuthor(db)
	book := BookData(db, author)

	t.Run("Accept Book Details Fetch", func(t *testing.T) {
		url := fmt.Sprintf("%s/book-detail/%s", baseUrl, book.Slug)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book details fetched successfully", body["message"])
	})
}

func adminGetBookContracts(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string, token string) {
	author := TestAuthor(db)
	BookData(db, author)

	t.Run("Reject Book Contracts Fetch Due To Invalid Contract Status", func(t *testing.T) {
		url := fmt.Sprintf("%s/contracts?contract_status=invalid-contract", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid contract status", body["message"])
	})

	t.Run("Accept Book Contracts Fetch", func(t *testing.T) {
		url := fmt.Sprintf("%s/contracts", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book Contracts fetched successfully", body["message"])
	})
}

func TestAdminBooks(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	admin := TestAdmin(db)
	token := AccessToken(db, admin)
	baseUrl := "/api/v1/admin/books"

	// Run Admin Books Endpoint Tests
	addBookGenre(t, app, db, baseUrl, token)
	addBookTag(t, app, db, baseUrl, token)
	updateBookGenre(t, app, db, baseUrl, token)
	updateBookTag(t, app, db, baseUrl, token)
	deleteBookGenre(t, app, db, baseUrl, token)
	deleteBookTag(t, app, db, baseUrl, token)
	adminGetBooks(t, app, db, baseUrl, token)
	adminGetBooksByAuthor(t, app, db, baseUrl, token)
	adminGetBookDetails(t, app, db, baseUrl, token)
	adminGetBookContracts(t, app, db, baseUrl, token)
}
