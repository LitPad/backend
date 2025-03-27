package database

import (
	"fmt"
	"log"
	"os"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Models() []interface{} {
	return []interface{}{
		// general
		&models.SiteDetail{},
		&models.Subscriber{},

		// accounts
		&models.User{},
		&models.AuthToken{},

		// book
		&models.Tag{},
		&models.Genre{},
		&models.Book{},
		&models.Chapter{},
		&models.Gift{},
		&models.SentGift{},
		&models.Comment{},
		&models.Reply{},
		&models.Vote{},

		// wallet
		&models.Coin{},
		&models.SubscriptionPlan{},
		&models.Transaction{},

		// waitlist
		&models.Waitlist{},

		// accounts
		&models.Notification{},

		// logs
		&models.Log{},
	}
}

func MakeMigrations(db *gorm.DB) {
	modelsList := Models()
	for _, model := range modelsList {
		db.AutoMigrate(model)
	}
}

func CreateTables(db *gorm.DB) {
	modelsList := Models()
	for _, model := range modelsList {
		db.Migrator().CreateTable(model)
	}
}

func DropTables(db *gorm.DB) {
	models := Models()
	for _, model := range models {
		db.Migrator().DropTable(model)
	}
}

func ConnectDb(cfg config.Config, loggedOpts ...bool) *gorm.DB {
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s"
	dbName := cfg.PostgresDB
	if os.Getenv("ENVIRONMENT") == "test" {
		dbName = cfg.TestPostgresDB
	}
	dsn := fmt.Sprintf(
		dsnTemplate,
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		dbName,
		cfg.PostgresPort,
		"UTC",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}
	log.Println("Connected to the database successfully")
	if len(loggedOpts) == 0 {
		db.Logger = logger.Default.LogMode(logger.Info)
	} else {
		db.Logger = logger.Default.LogMode(logger.Error)
	}

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		log.Fatal("failed to create extension: " + result.Error.Error())
	}

	// Add Migrations
	if os.Getenv("ENVIRONMENT") != "test" {
		log.Println("Running Migrations")
		MakeMigrations(db)
	}
	return db
}
