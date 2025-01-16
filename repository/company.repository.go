package repository

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	CreateCompany(company *model.Company) error
	GetUserCompanies(userEmail string) ([]model.Company, error)
	GetUserCompanyByName(email string, companyName string) (*model.Company, error)
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

func (cr *companyRepository) GetUserCompanies(userEmail string) ([]model.Company, error) {
	var companies []model.Company
	if err := cr.db.Where("owner_email = ?", userEmail).Find(&companies).Error; err != nil {
		return nil, err
	}

	return companies, nil
}

func (cr *companyRepository) GetUserCompanyByName(email string, companyName string) (*model.Company, error) {
	var company model.Company
	result := cr.db.Where("owner_email = ? and name = ?", email, companyName).First(&company)
	if result.Error != nil {
		return nil, result.Error
	}

	return &company, nil
}
