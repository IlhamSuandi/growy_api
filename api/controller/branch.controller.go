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

type BranchController struct {
	Logger         *logrus.Logger
	BranchUsecase  usecase.BranchUsecase
	CompanyUsecase usecase.CompanyUsecase
}

func NewBranchController(branchUsecase usecase.BranchUsecase, companyUsecase usecase.CompanyUsecase) *BranchController {
	return &BranchController{
		Logger:         utils.Log,
		BranchUsecase:  branchUsecase,
		CompanyUsecase: companyUsecase,
	}
}

// @Tags Branch
// @Summary [owner only] get company branches
// @Description getting company branches that only owner can get [owner only]
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company_name query string true "company name"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 201 {object} types.Response{data=model.Branch} "Successfully created branch"
// @Router /branch [get]
func (bc *BranchController) GetCompanyBranches(w http.ResponseWriter, r *http.Request) {
	bc.Logger.Info("[/branch] getting company name from query params")
	companyName := r.URL.Query().Get("company_name")

	if companyName == "" {
		bc.Logger.Errorf("[/branch] company name is required")
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "company name is required",
			Error:   "company name is required",
			Status:  http.StatusBadRequest,
		})
		return
	}

	bc.Logger.Info("[/branch] getting user informations from auth middleware")
	userInfo := r.Context().Value("userInfo").(*model.User)

	bc.Logger.Info("[/branch] getting user company")
	userCompany, err := bc.CompanyUsecase.GetUserCompanyByName(userInfo.Email, companyName)
	if err != nil {
		bc.Logger.Errorf("[/branch] error getting user company %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error getting user company",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	bc.Logger.Info("[/branch] getting company branches")
	userBranches, err := bc.BranchUsecase.GetCompanyBranches(userCompany.Id)
	if err != nil {
		bc.Logger.Errorf("[/branch] error getting company branches %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error getting company branches",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "successfully getting company branches",
		Data:    userBranches,
		Status:  http.StatusOK,
	})
}

// @Tags Branch
// @Summary [owner only] create new branch using company name
// @Description creating new branch using company name that only owner can create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateBranchRequest true "Request body"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 201 {object} types.Response{data=model.Branch} "Successfully created branch"
// @Router /branch [post]
func (bc *BranchController) CreateBranch(w http.ResponseWriter, r *http.Request) {
	bc.Logger.Info("[/branch] creating branch")
	bc.Logger.Info("[/branch] parsing request body")
	var payload dto.CreateBranchRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		bc.Logger.Errorf("[/branch] error parsing request body %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "error parsing json",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	bc.Logger.Info("[/branch] getting user info from auth middleware")
	userInfo := r.Context().Value("userInfo").(*model.User)
	userCompany, err := bc.CompanyUsecase.GetUserCompanyByName(userInfo.Email, payload.CompanyName)
	if err != nil {
		bc.Logger.Errorf("[/branch] error getting user company %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error getting user company",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	bc.Logger.Info("[/branch] creating branch")
	branch := model.Branch{
		CompanyId: userCompany.Id,
		Name:      payload.BranchName,
		Address:   payload.Address,
	}

	if err := bc.BranchUsecase.CreateBranch(&branch); err != nil {
		bc.Logger.Errorf("[/branch] error creating branch %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error creating branch",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	bc.Logger.Info("[/branch] successfully creating branch")
	response.WriteJSON(w, http.StatusCreated, types.Response{
		Message: "Successfully created branch",
		Data:    branch,
		Status:  http.StatusCreated,
	})
}
