package repository

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type BranchRepository interface {
	GetCompanyBranchByName(companyId uint, branchName string) (*model.Branch, error)
	GetCompanyBranches(companyId uint) ([]model.Branch, error)
	CreateBranch(branch *model.Branch) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{
		db: db,
	}
}

func (br *branchRepository) GetCompanyBranchByName(companyId uint, branchName string) (*model.Branch, error) {
	var branch model.Branch
	result := br.db.Where("company_id = ? and name = ?", companyId, branchName).First(&branch)
	if result.Error != nil {
		return nil, result.Error
	}

	return &branch, nil
}

func (br *branchRepository) GetCompanyBranches(companyId uint) ([]model.Branch, error) {
	var branches []model.Branch
	result := br.db.Where("company_id = ?", companyId).Preload("Employees.Salary").Find(&branches)
	if result.Error != nil {
		return nil, result.Error
	}

	return branches, nil
}

func (br *branchRepository) CreateBranch(branch *model.Branch) error {
	result := br.db.Create(branch)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("error creating branch")
	}

	return nil
}
