package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/LitPad/backend/database"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getBookTags(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Accept Book Tags Fetch", func(t *testing.T) {
		TagData(db) // Get or create tag
		url := fmt.Sprintf("%s/tags", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		res := ProcessTestGetOrDelete(app, url, "GET")
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid book genre", body["message"])
	})

	t.Run("Accept Books Fetch", func(t *testing.T) {
		res := ProcessTestGetOrDelete(app, baseUrl, "GET")
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
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		url := fmt.Sprintf("%s/author/%s", baseUrl, user.Username)
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		url := fmt.Sprintf("%s/book/%s/chapters", baseUrl, book.Slug)
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		res := ProcessTestGetOrDelete(app, url, "GET")
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
		url := fmt.Sprintf("%s/book/%s", baseUrl, book.Slug)
		res := ProcessTestGetOrDelete(app, url, "GET")
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book details fetched successfully", body["message"])
	})
}

func createBook(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	bookData := schemas.BookCreateSchema{
		Title: "Test Book Title", Blurb: "Test Blurb",
		GenreSlug: "invalid", TagSlugs: []string{"slug"}, AgeDiscretion: choices.ATYPE_EIGHTEEN,
	}
	author := TestAuthor(db)
	token := AccessToken(db, author)

	t.Run("Reject Book Creation Due To Invalid Genre Slug", func(t *testing.T) {
		res := ProcessMultipartTestBody(t, app, baseUrl, "POST", bookData, "", "", token)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Invalid genre slug!", body["data"].(map[string]interface{})["genre_slug"])
	})
	genre := GenreData(db)
	t.Run("Reject Book Creation Due To Invalid Tag Slugs", func(t *testing.T) {
		bookData.GenreSlug = genre.Slug
		bookData.TagSlugs = []string{"invalid"}
		res := ProcessMultipartTestBody(t, app, baseUrl, "POST", bookData, "", "", token)

		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "The following are invalid tag slugs: invalid", body["data"].(map[string]interface{})["tag_slugs"])
	})

	t.Run("Reject Book Creation Due To Empty File", func(t *testing.T) {
		bookData.TagSlugs = []string{genre.Tags[0].Slug}
		res := ProcessMultipartTestBody(t, app, baseUrl, "POST", bookData, "", "", token)

		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Image is required", body["data"].(map[string]interface{})["cover_image"])
	})

	t.Run("Accept Book Creation Due To Valid Data", func(t *testing.T) {
		// Create a temporary file
		tempFilePath := CreateTempImageFile(t)
		defer os.Remove(tempFilePath)
		res := ProcessMultipartTestBody(t, app, baseUrl, "POST", bookData, "cover_image", tempFilePath, token)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book created successfully", body["message"])
	})
}

func updateBook(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	genre := GenreData(db)
	bookData := schemas.BookCreateSchema{
		Title: "Test Book Title Updated", Blurb: "Test Blurb Updated",
		GenreSlug: genre.Slug, TagSlugs: []string{genre.Tags[0].Slug}, AgeDiscretion: choices.ATYPE_EIGHTEEN,
	}
	author := TestAuthor(db)
	token := AccessToken(db, author)
	book := BookData(db, author)

	t.Run("Reject Book Update Due To Invalid Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/invalid-slug", baseUrl)
		res := ProcessMultipartTestBody(t, app, url, "PUT", bookData, "", "", token)
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Author has no book with that slug", body["message"])
	})

	t.Run("Accept Book Update Due To Valid Data", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/%s", baseUrl, book.Slug)
		res := ProcessMultipartTestBody(t, app, url, "PUT", bookData, "", "", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book updated successfully", body["message"])
	})
}

func deleteBook(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	author := TestAuthor(db)
	token := AccessToken(db, author)
	book := BookData(db, author)

	t.Run("Accept Book Delete Due To Valid Data", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/%s", baseUrl, book.Slug)
		res := ProcessTestGetOrDelete(app, url, "DELETE", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Book deleted successfully", body["message"])
	})
}

func addChapter(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	author := TestAuthor(db)
	token := AccessToken(db, author)
	book := BookData(db, author)
	chapterData := schemas.ChapterCreateSchema{
		Title: "Test Chapter Title", Text: "Test Content",
	}
	t.Run("Accept Chapter Creation Due To Valid Data", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/%s/add-chapter", baseUrl, book.Slug)
		res := ProcessJsonTestBody(t, app, url, "POST", chapterData, token)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Chapter added successfully", body["message"])
	})
}

func updateChapter(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	author := TestAuthor(db)
	token := AccessToken(db, author)

	invalidOwner := TestAuthor(db, true)
	invalidOwnerToken := AccessToken(db, invalidOwner)
	book := BookData(db, author)
	chapter := ChapterData(db, book)

	chapterData := schemas.ChapterCreateSchema{
		Title: "Test Chapter Title Updated", Text: "Test Content Updated",
	}

	t.Run("Reject Chapter Update Due To Invalid Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/chapter/invalid-slug", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "PUT", chapterData, token)
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "No chapter with that slug", body["message"])
	})

	t.Run("Reject Chapter Update Due To Invalid Owner", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/chapter/%s", baseUrl, chapter.Slug)
		res := ProcessJsonTestBody(t, app, url, "PUT", chapterData, invalidOwnerToken)
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Not yours to edit", body["message"])
	})

	t.Run("Accept Chapter Update Due To Valid Data", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/chapter/%s", baseUrl, chapter.Slug)
		res := ProcessJsonTestBody(t, app, url, "PUT", chapterData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Chapter updated successfully", body["message"])
	})
}

func deleteChapter(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	author := TestAuthor(db)
	token := AccessToken(db, author)

	book := BookData(db, author)
	chapter := ChapterData(db, book)

	t.Run("Accept Chapter Delete Due To Valid Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s/book/chapter/%s", baseUrl, chapter.Slug)
		res := ProcessTestGetOrDelete(app, url, "DELETE", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Chapter deleted successfully", body["message"])
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
	createBook(t, app, db, baseUrl)
	updateBook(t, app, db, baseUrl)
	deleteBook(t, app, db, baseUrl)
	addChapter(t, app, db, baseUrl)
	updateChapter(t, app, db, baseUrl)
	deleteChapter(t, app, db, baseUrl)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
