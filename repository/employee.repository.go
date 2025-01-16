package repository

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetUserEmployee(email string) (*model.Employee, error)
	CreateEmployee(employee *model.Employee) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (er *employeeRepository) GetUserEmployee(email string) (*model.Employee, error) {
	var employee model.Employee
	result := er.db.Where("employee_email = ?", email).First(&employee)
	if result.Error != nil {
		return nil, result.Error
	}

	return &employee, nil
}

func (er *employeeRepository) CreateEmployee(employee *model.Employee) error {
	result := er.db.Create(&employee)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("error creating employee")
	}

	return nil
}
