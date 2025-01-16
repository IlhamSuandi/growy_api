package usecase

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type CompanyUsecase interface {
	CreateCompany(company *model.Company) error
	GetUserCompanies(email string) ([]model.Company, error)
  GetUserCompanyByName(email string, companyName string) (*model.Company, error)
}

type companyUsecase struct {
	logger      *logrus.Logger
	CompanyRepo repository.CompanyRepository
}

func NewCompanyUsecase(companyRepo repository.CompanyRepository) CompanyUsecase {
	return &companyUsecase{
		logger:      utils.Log,
		CompanyRepo: companyRepo,
	}
}

func (cu *companyUsecase) CreateCompany(company *model.Company) error {
	return cu.CompanyRepo.CreateCompany(company)
}

func (cu *companyUsecase) GetUserCompanies(email string) ([]model.Company, error) {
	return cu.CompanyRepo.GetUserCompanies(email)
}

func (cu *companyUsecase) GetUserCompanyByName(email string, companyName string) (*model.Company, error) {
  return cu.CompanyRepo.GetUserCompanyByName(email, companyName)
}
