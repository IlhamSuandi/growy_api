package usecase

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type SalaryUsecase interface{
  CreateEmployeeSalary(salary *model.Salary) error
}

type salaryUsecase struct {
	Logger           *logrus.Logger
	SalaryRepository repository.SalaryRepository
}

func NewSalaryUsecase(salaryRepo repository.SalaryRepository) SalaryUsecase {
	return &salaryUsecase{
		Logger:           utils.Log,
		SalaryRepository: salaryRepo,
	}
}

func (sc *salaryUsecase) CreateEmployeeSalary(salary *model.Salary) error {
	return sc.SalaryRepository.CreateEmployeeSalary(salary)
}
