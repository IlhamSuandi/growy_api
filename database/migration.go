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
		&model.User{},
		&model.Log{},
		&model.Session{},
		&model.Permission{},
		&model.Attendance{},
		&model.QRCode{},
		&model.Role{},
		&model.Company{},
		&model.Branch{},
		&model.CompanyOption{},
		&model.Employee{},
		&model.WorkSchedule{},
		&model.Salary{},
	); err != nil {
		log.Error(err)
	}

	defer log.Info("successfully migrated")
}
