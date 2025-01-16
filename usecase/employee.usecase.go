package usecase

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type EmployeeUsecase interface {
	GetUserEmployee(email string) (*model.Employee, error)
	CreateEmployee(employee *model.Employee) error
}

type employeeUsecase struct {
	Logger             *logrus.Logger
	EmployeeRepository repository.EmployeeRepository
}

func NewEmployeeUsecase(employeeRepository repository.EmployeeRepository) EmployeeUsecase {
	return &employeeUsecase{
		Logger:             utils.Log,
		EmployeeRepository: employeeRepository,
	}
}

func (eu *employeeUsecase) GetUserEmployee(email string) (*model.Employee, error) {
	return eu.EmployeeRepository.GetUserEmployee(email)
}

func (eu *employeeUsecase) CreateEmployee(employee *model.Employee) error {
	return eu.EmployeeRepository.CreateEmployee(employee)
}
