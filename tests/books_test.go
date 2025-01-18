package tests

import (
	"fmt"
	"testing"

	"github.com/LitPad/backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getBookTags(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Accept Book Tags Fetch", func(t *testing.T) {
		TagData(db) // Get or create tag
		url := fmt.Sprintf("%s/tags", baseUrl)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tags fetched successfully", body["message"])
	})
}

func getBookGenres(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Accept Book Genres Fetch", func(t *testing.T) {
		GenreData(db) // Get or create tag
		url := fmt.Sprintf("%s/genres", baseUrl)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Genres fetched successfully", body["message"])
	})
}

func getBooks(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := TestAuthor(db)
	BookData(db, user) // Get or create book
	t.Run("Reject Books Fetch Due To Invalid Genre Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s?genre_slug=invalid-genre", baseUrl)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid book genre", body["message"])
	})

	t.Run("Accept Books Fetch", func(t *testing.T) {
		res := ProcessTestGet(app, baseUrl)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Books fetched successfully", body["message"])
	})
}

func getBooksByAuthor(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject Books Fetch Due To Invalid Username", func(t *testing.T) {
		url := fmt.Sprintf("%s/author/invalid-username", baseUrl)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid author username", body["message"])
	})

	t.Run("Accept Author Books Fetch", func(t *testing.T) {
		user := TestAuthor(db)
		BookData(db, user) // Get or create book
		t.Log(user)
		url := fmt.Sprintf("%s/author/%s", baseUrl, user.Username)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Books fetched successfully", body["message"])
	})
}

func getBookChapters(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject Book Chapters Fetch Due To Invalid Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/invalid-slug/chapters", baseUrl)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "No book with that slug", body["message"])
	})

	t.Run("Accept Book Chapters Fetch", func(t *testing.T) {
		user := TestAuthor(db)
		book := BookData(db, user) // Get or create book
		ChapterData(db, book)
		t.Log(user)
		url := fmt.Sprintf("%s/book/%s/chapters", baseUrl, book.Slug)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Chapters fetched successfully", body["message"])
	})
}

func getBook(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject Book Details Fetch Due To Invalid Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/invalid-slug", baseUrl)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "No book with that slug", body["message"])
	})

	t.Run("Accept Book Details Fetch", func(t *testing.T) {
		user := TestAuthor(db)
		book := BookData(db, user) // Get or create book
		ChapterData(db, book)
		t.Log(user)
		url := fmt.Sprintf("%s/book/%s", baseUrl, book.Slug)
		res := ProcessTestGet(app, url)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book details fetched successfully", body["message"])
	})
}

func TestBooks(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	baseUrl := "/api/v1/books"

	// Run Book Endpoint Tests
	getBookTags(t, app, db, baseUrl)
	getBookGenres(t, app, db, baseUrl)
	getBooks(t, app, db, baseUrl)
	getBooksByAuthor(t, app, db, baseUrl)
	getBookChapters(t, app, db, baseUrl)
	getBook(t, app, db, baseUrl)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
