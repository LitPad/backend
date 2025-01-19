package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/database"
	"github.com/LitPad/backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func CreateSingleTable(db *gorm.DB, model interface{}) {
	db.AutoMigrate(&model)
}

func DropAndCreateSingleTable(db *gorm.DB, model interface{}) {
	db.Migrator().DropTable(&model)
	db.AutoMigrate(&model)
}

func SetupTestDatabase(t *testing.T) *gorm.DB {
	cfg := config.GetConfig()
	return database.ConnectDb(cfg, false)
}

func CloseTestDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection: " + err.Error())
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal("Failed to close database connection: " + err.Error())
	}
}


func Setup(t *testing.T, app *fiber.App) *gorm.DB {
	t.Setenv("ENVIRONMENT", "test")
	t.Setenv("CONFIG_PATH", "../")

	// Set up the test database
	db := SetupTestDatabase(t)

	routes.SetupRoutes(app, db)
	t.Logf("Making Database Migrations....")
	database.DropTables(db)
	database.CreateTables(db)
	t.Logf("Database Migrations Made successfully")
	return db
}

func ParseResponseBody(t *testing.T, b io.ReadCloser) interface{} {
	body, _ := io.ReadAll(b)
	// Parse the response body as JSON
	responseBody := make(map[string]interface{})
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Errorf("error parsing response body as JSON: %s", err)
	}
	return responseBody
}

func ProcessTestGet(app *fiber.App, url string, access ...string) *http.Response {
	req := httptest.NewRequest("GET", url, nil)
	if access != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access[0]))
	}
	res, _ := app.Test(req)
	return res
}

func ProcessTestBody(t *testing.T, app *fiber.App, url string, method string, body interface{}, access ...string) *http.Response {
	// Marshal the test data to JSON
	requestBytes, err := json.Marshal(body)
	requestBody := bytes.NewReader(requestBytes)
	assert.Nil(t, err)
	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	if access != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access[0]))
	}
	res, err := app.Test(req)
	if err != nil {
		log.Println(err)
	}
	return res
}

func ProcessMultipartTestBody(t *testing.T, app *fiber.App, url string, method string, body interface{}, fileFieldName string, filePath string, access ...string) *http.Response {
	var requestBody *bytes.Buffer
	// Multipart handling
	requestBody = &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	// Populate multipart form fields and files from the struct
	populateMultipartFromStruct(t, body, writer)

	// Add the file separately if provided
	if filePath != "" {
		file, err := os.Open(filePath)
		assert.Nil(t, err)
		defer file.Close()

		// Add the file to the multipart form
		fileWriter, err := writer.CreateFormFile(fileFieldName, filepath.Base(filePath))
		assert.Nil(t, err)
		_, err = io.Copy(fileWriter, file)
		assert.Nil(t, err)
	}

	// Close the writer to finalize the body
	err := writer.Close()
	assert.Nil(t, err)

	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if len(access) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access[0]))
	}

	res, err := app.Test(req, 10000)
	assert.Nil(t, err)
	return res
}

// Helper: Populate multipart writer from struct fields
func populateMultipartFromStruct(t *testing.T, body interface{}, writer *multipart.Writer) {
	val := reflect.ValueOf(body)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Tag.Get("form")
		if fieldName == "" {
			fieldName = fieldType.Name // Use field name if no tag is specified
		}

		// Handle supported types
		switch field.Kind() {
		case reflect.String:
			// Add string fields
			err := writer.WriteField(fieldName, field.String())
			assert.Nil(t, err)

		case reflect.Slice:
			// Handle slices (e.g., []string)
			if field.Type().Elem().Kind() == reflect.String {
				for i := 0; i < field.Len(); i++ {
					err := writer.WriteField(fieldName, field.Index(i).String())
					assert.Nil(t, err)
				}
			} else {
				t.Errorf("Unsupported slice element type for field: %s", fieldName)
			}

		case reflect.Ptr:
			// Handle file pointers
			if file, ok := field.Interface().(*os.File); ok && file != nil {
				fileWriter, err := writer.CreateFormFile(fieldName, filepath.Base(file.Name()))
				assert.Nil(t, err)
				_, err = io.Copy(fileWriter, file)
				assert.Nil(t, err)
				file.Close()
			} else {
				t.Errorf("Unsupported pointer type for field: %s", fieldName)
			}

		default:
			// Handle all other types generically
			if field.CanInterface() {
				value := field.Interface()

				// Check if the type implements fmt.Stringer
				if stringer, ok := value.(fmt.Stringer); ok {
					err := writer.WriteField(fieldName, stringer.String())
					assert.Nil(t, err)
				} else if reflect.TypeOf(value).ConvertibleTo(reflect.TypeOf("")) {
					// Attempt to convert the value to a string
					err := writer.WriteField(fieldName, fmt.Sprintf("%v", value))
					assert.Nil(t, err)
				} else {
					t.Errorf("Unsupported field type: %s (Field: %s)", field.Type(), fieldName)
				}
			}
		}
	}
}

func CreateTempImageFile(t *testing.T) string {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test-image-*.jpg")
	assert.Nil(t, err)

	// Close the file so we can write to it
	defer tempFile.Close()

	// Write dummy image data to the file (a minimal valid JPEG header)
	_, err = tempFile.Write([]byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00\x01\x01\x01\x00H\x00H\x00\x00"))
	assert.Nil(t, err)

	// Return the path of the temporary file
	return tempFile.Name()
}

func RemoveCreatedAndUpdated (body map[string]interface{}, dataType string) {
	// To remove created_at and updated_at
	dataMap := body["data"].(map[string]interface{})
	dataMapValues := dataMap[dataType].([]interface{})
	dataMapValue := dataMapValues[0].(map[string]interface{})
	delete(dataMapValue, "created_at")
	delete(dataMapValue, "updated_at")
}

func ConvertDateTime(timeObj time.Time) string {
	roundedTime := timeObj.Round(time.Microsecond)
	formatted := roundedTime.Format("2006-01-02T15:04:05")

	// Get the microsecond part and round it
	microseconds := roundedTime.Nanosecond() / 1000

	// Append the rounded microsecond part to the formatted string
	formatted = fmt.Sprintf("%s.%06d", formatted, microseconds)
	formatted = strings.TrimRight(formatted, "0")
	// Append the timezone information
	formatted = fmt.Sprintf("%s%s", formatted, roundedTime.Format("-07:00"))

	return formatted
}