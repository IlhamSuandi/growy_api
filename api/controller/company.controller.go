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

type CompanyController struct {
	Logger         *logrus.Logger
	CompanyUsecase usecase.CompanyUsecase
}

func NewCompanyController(companyUsecase usecase.CompanyUsecase) *CompanyController {
	return &CompanyController{
		Logger:         utils.Log,
		CompanyUsecase: companyUsecase,
	}
}

func (cc *CompanyController) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateCompanyRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		cc.Logger.Errorf("[%s /company] error parsing request body %s", r.Method, err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error request parsing body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	cc.Logger.Infof("[%s /company] getting user informations from auth middleware", r.Method)
	userInfo := r.Context().Value("userInfo").(*model.User)

	cc.Logger.Infof("[%s /company] creating company", r.Method)
	company := model.Company{
		OwnerEmail: userInfo.Email,
		Name:       payload.Name,
		Address:    payload.Address,
	}

	if err := cc.CompanyUsecase.CreateCompany(&company); err != nil {
		cc.Logger.Errorf("[%s /company] error creating company %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully created company",
		Data: dto.CreateCompanyResponse{
			Name:       company.Name,
			Address:    company.Address,
			OwnerEmail: company.OwnerEmail,
		},
		Status: http.StatusOK,
	})
}

// @Tags Company
// @Summary get owner companies
// @Description get all companies of owner
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 200 {object} types.Response{data=dto.GetUserCompaniesResponse} "Successfully get owner companies"
// @Router /company [get]
func (cc *CompanyController) GetUserCompanies(w http.ResponseWriter, r *http.Request) {
	cc.Logger.Infof("[%s /company] getting user informations", r.Method)
	userInfo := r.Context().Value("userInfo").(*model.User)

	cc.Logger.Infof("[%s /company] getting user companies", r.Method)
	companies, err := cc.CompanyUsecase.GetUserCompanies(userInfo.Email)
	if err != nil {
		cc.Logger.Errorf("[%s /company] error getting user companies %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully get user companies",
		Data: dto.GetUserCompaniesResponse{
			Company: companies,
		},
		Status: http.StatusOK,
	})
}
