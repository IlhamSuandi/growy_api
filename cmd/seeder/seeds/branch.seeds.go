package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

var MainBranch model.Branch

func SeedBranch(db *gorm.DB) {
	log := utils.Log
	log.Info("seeding branch")

	MainBranch = model.Branch{
		CompanyId: Company.Id,
		Name:      "main",
		Address:   "jakarta",
	}

	result := db.Save(&MainBranch)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	if result.RowsAffected == 0 {
		log.Fatal("no rows affected")
	}

	UserOwner.IsOnBoarded = true

	updateResult := db.Save(&UserOwner)
	if updateResult.Error != nil {
		log.Fatal(result.Error)
	}

	if updateResult.RowsAffected == 0 {
		log.Fatal("no rows affected")
	}

	log.Info("successfully seeded branch")
}
