package repository

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type SalaryRepository interface {
	CreateEmployeeSalary(salary *model.Salary) error
}

type salaryRepository struct {
	db *gorm.DB
}

func NewSalaryRepository(db *gorm.DB) SalaryRepository {
	return &salaryRepository{
		db: db,
	}
}

func (sr *salaryRepository) CreateEmployeeSalary(salary *model.Salary) error {
	result := sr.db.Create(&salary)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("error creating salary")
	}

	return nil
}
