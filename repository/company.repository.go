package repository

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	CreateCompany(company *model.Company) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

func (cr *companyRepository) CreateCompany(company *model.Company) error {
	createdCompany := cr.db.Create(&company)
	if createdCompany.Error != nil {
		return createdCompany.Error
	}

	if createdCompany.RowsAffected == 0 {
		return errors.New("error creating company")
	}

	return nil
}
