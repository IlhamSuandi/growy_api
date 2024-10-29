package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

var Company model.Company

func SeedCompany(db *gorm.DB) error {
	log := utils.Log
	log.Info("seeding company")

	Company = model.Company{
		Name:       "growy",
		Address:    "jakarta",
		OwnerEmail: UserAdmin.Email,
	}

	companyRepo := repository.NewCompanyRepository(db)
	if err := companyRepo.CreateCompany(&Company); err != nil {
		return err
	}

	log.Info("successfully seeded company")
	return nil
}
