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
		&models.Token{},

		// book
		&models.Tag{},
		&models.Genre{},
		&models.Book{},
		&models.Chapter{},
		&models.Gift{},
		&models.SentGift{},
		&models.Review{},
		&models.Reply{},

		// wallet
		&models.Coin{},
		&models.Transaction{},
		&models.BoughtBook{},
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
	// Drop Tables
	models := Models()
	for _, model := range models {
		db.Migrator().DropTable(model)
	}
}

func ConnectDb(cfg config.Config) *gorm.DB {
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s"
	dsn := fmt.Sprintf(
		dsnTemplate,
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
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
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		log.Fatal("failed to create extension: " + result.Error.Error())
	}

	// Add Migrations
	MakeMigrations(db)
	return db
}
