package seeds

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

var MainBranch model.Branch

func SeedBranch(db *gorm.DB) error {
	log := utils.Log
	log.Info("seeding branch")

	employees := []model.User{
		UserOne,
		UserTwo,
		UserThree,
		UserFour,
		UserFive,
	}
	MainBranch = model.Branch{
		CompanyId: Company.Id,
		Name:      "main",
		Address:   "jakarta",
		Employees: employees,
	}

	result := db.Save(&MainBranch)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	log.Info("successfully seeded branch")
	return nil
}
