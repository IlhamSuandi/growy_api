package database

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	log := utils.Log
	log.Info("Auto Migrating Database...")

	if err := db.AutoMigrate(
		&model.Log{},
		&model.User{},
		&model.Session{},
		&model.Permission{},
		&model.Attendance{},
		&model.QRCode{},
		&model.Role{},
		&model.Company{},
		&model.Branch{},
		&model.BranchOption{},
	); err != nil {
		log.Println(err)
	}

	defer log.Info("successfully migrated")
}
