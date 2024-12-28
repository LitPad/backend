package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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
	t.Setenv("ENVIRONMENT", "TESTING")
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

func RemoveCreatedAndUpdated (body map[string]interface{}, dataType string) {
	// To remove created_at and updated_at
	dataMap := body["data"].(map[string]interface{})
	dataMapValues := dataMap[dataType].([]interface{})
	dataMapValue := dataMapValues[0].(map[string]interface{})
	delete(dataMapValue, "created_at")
	delete(dataMapValue, "updated_at")
}