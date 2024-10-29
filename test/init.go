package test

import (
	"github.com/ilhamSuandi/business_assistant/database"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	utils.Log.Info("Initializing test environment")

	db, err := database.Connect("localhost", "test_db")
	if err != nil {
		panic("failed to connect to database")
	}

	database.AutoMigrate(db)

	DB = db

	utils.RegisterValidator()
}
