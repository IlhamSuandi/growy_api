package database

import (
	"fmt"

	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(host string, dbname string) (*gorm.DB, error) {
	user := config.DB_USER
	password := config.DB_PASSWORD
	port := config.DB_PORT
	sslmode := config.DB_SSLMODE
	timezone := config.DB_TIMEZONE

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
    PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	defer utils.Log.Info("Database connection established")
	return db, nil
}

func CloseDb(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		utils.Log.Errorf("error getting database instance %s", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("error closing database instance %s", err)
		return
	}
	utils.Log.Info("Database connection closed")
}
