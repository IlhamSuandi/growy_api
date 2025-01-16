package controller

import (
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type SalaryController struct {
	Logger        *logrus.Logger
	SalaryUsecase usecase.SalaryUsecase
}

func NewSalaryController(salaryUsecase usecase.SalaryUsecase) *SalaryController {
	return &SalaryController{
		Logger:        utils.Log,
		SalaryUsecase: salaryUsecase,
	}
}

func (sc *SalaryController) CreateEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateSalaryRequest

	sc.Logger.Info("[/salary] parsing request body")
	if err := utils.ParseJSON(r, &payload); err != nil {
		sc.Logger.Errorf("[/salary] error parsing request body %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error Parsing Request Body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	sc.Logger.Info("[/salary] creating salary")
	salary := model.Salary{
		EmployeeEmail: payload.EmployeeEmail,
	}

	if payload.Monthly == nil && payload.Hourly == nil {
		sc.Logger.Error("[/salary] Both Monthly and Hourly salary are nil")
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Monthly or Hourly salary is required",
			Error:   "Both Monthly and Hourly salary are nil",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if payload.Monthly != nil {
		sc.Logger.Info("[/salary] creating monthly salary")
		salary.Monthly = *payload.Monthly
	}

	if payload.Hourly != nil {
		sc.Logger.Info("[/salary] creating hourly salary")
		salary.Hourly = *payload.Hourly
	}

	if err := sc.SalaryUsecase.CreateEmployeeSalary(&salary); err != nil {
		sc.Logger.Errorf("[/salary] error creating salary %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error creating salary",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

  response.WriteJSON(w, http.StatusCreated, types.Response{
    Message: "Salary created successfully",
    Data: dto.CreateSalaryResponse{
      EmployeeEmail: salary.EmployeeEmail,
      Monthly:       salary.Monthly,
      Hourly:        salary.Hourly,
    },
    Status:  http.StatusCreated,
  })
}
