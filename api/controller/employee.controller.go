package controller

import (
	"fmt"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/auth"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type EmployeeController struct {
	Logger          *logrus.Logger
	EmailUsecase    usecase.EmailUsecase
	EmployeeUsecase usecase.EmployeeUsecase
	UserUsecase     usecase.UserUsecase
	companyUsecase  usecase.CompanyUsecase
	BranchUsecase   usecase.BranchUsecase
}

func NewEmployeeController(
	employeeUsecase usecase.EmployeeUsecase,
	userUsecase usecase.UserUsecase,
	companyUsecase usecase.CompanyUsecase,
	branchUsecase usecase.BranchUsecase,
) *EmployeeController {
	return &EmployeeController{
		Logger:          utils.Log,
		EmailUsecase:    usecase.NewEmailUsecase(),
		EmployeeUsecase: employeeUsecase,
		UserUsecase:     userUsecase,
		companyUsecase:  companyUsecase,
		BranchUsecase:   branchUsecase,
	}
}

// TODO: fix this
// @Tags Employee
// @Summary Add Employee
// @Description invite user to join company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddEmployeeRequest true "Request body"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 201 {object} types.Response{data=dto.AddEmployeeResponse} "Successfully Added Employee"
// @Success 200 {object} types.Response{data=dto.AddEmployeeResponse} "Successfully sent invitation email"
// @Router /employee [post]
func (ec *EmployeeController) AddEmployee(w http.ResponseWriter, r *http.Request) {
	var payload dto.AddEmployeeRequest

	ec.Logger.Info("[/employee] parsing request body")
	if err := utils.ParseJSON(r, &payload); err != nil {
		ec.Logger.Errorf("[/employee] error parsing request body %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "error parsing json",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	ec.Logger.Info("[/employee] getting user info from auth middleware")
	userInfo := r.Context().Value("userInfo").(*model.User)

	ec.Logger.Info("[/employee] getting user company")
	userCompany, err := ec.companyUsecase.GetUserCompanyByName(userInfo.Email, payload.CompanyName)
	if err != nil {
		ec.Logger.Errorf("[/employee] error getting user company %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error getting user company",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ec.Logger.Info("[/employee] getting user branch")
	userBranch, err := ec.BranchUsecase.GetCompanyBranchByName(userCompany.Id, payload.BranchName)
	if err != nil {
		ec.Logger.Errorf("[/employee] error getting user branch %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error getting user branch",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ec.Logger.Info("[/employee] checking if user is exists")
	existedUser, err := ec.UserUsecase.GetUserByEmail(payload.Email)
	if err == nil {
		ec.Logger.Info("[/employee] user exists")

		_, err := ec.EmployeeUsecase.GetUserEmployee(userInfo.Email)
		if err == nil {
			ec.Logger.Error("[/employee] user already registered as employee")
			response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
				Message: "Error Adding Employee",
				Error:   err,
				Status:  http.StatusBadRequest,
			})
			return
		}

		employee := model.Employee{
			EmployeeEmail: existedUser.Email,
			Position:      payload.Position,
			BranchId:      userBranch.Id,
			Pending:       true,
		}

		if err := ec.EmployeeUsecase.CreateEmployee(&employee); err != nil {
			ec.Logger.Errorf("[/employee] error creating employee %s", err)
			response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
				Message: "Error Adding Employee",
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			})
			return
		}

		response.WriteJSON(w, http.StatusCreated, types.Response{
			Message: "Successfully Added Employee",
			Data: dto.AddEmployeeResponse{
				Email:    existedUser.Email,
				Position: payload.Position,
			},
			Status: http.StatusCreated,
		})
		return
	}

	ec.Logger.Info("[/employee] user does not exist")
	ec.Logger.Info("[/employee] creating token")
	token, _, err := auth.CreateToken(payload.Email, userBranch.Id)
	if err != nil {
		ec.Logger.Errorf("[/employee] error creating token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Adding Employee",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	href := config.CLIENT_URL + "/accept-invitation/" + token
	subject := "Invitation to Join " + userCompany.Name
	body := fmt.Sprintf(`
    <h2 style="color: #555;">Invitation to Join %s</h2>
    <p>Dear %s,</p>
    <p>You have been invited to join <strong>%s</strong> as <strong>%s</strong>. We are thrilled to have the opportunity to welcome you to our team!</p>
    <p><strong>Here are the details of your invitation:</strong></p>
    <ul>
        <li><strong>Position:</strong> %s</li>
        <li><strong>Location:</strong> %s</li>
    </ul>
    <p>To accept this invitation and complete your registration, please click the link below:</p>
    <p><a href="%s" style="color: #1a73e8; text-decoration: none;">Accept Invitation</a></p>
    <p>If you have any questions or need assistance, please contact our support team at <a href="mailto:%s" style="color: #1a73e8;">%s</a>.</p>
    <p>Thank you, and we look forward to having you on board!</p>
    <p>Best regards,</p>
    <p><strong>%s</strong></p>
    <p style="font-size: 0.9em; color: #999;">*This is an automated message. Please do not reply to this email.*</p>
`, userCompany.Name, payload.Email, userCompany.Name, payload.Position, payload.Position, userBranch.Address, href, userInfo.Email, config.SMTPUSERNAME, userCompany.Name)

	ec.EmailUsecase.SendEmailAsync(payload.Email, subject, body)

	ec.Logger.Info("[/employee] Successfully sent invitation email")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully sent invitation email to " + payload.Email,
		Data: dto.AddEmployeeResponse{
			Email:    payload.Email,
			Position: payload.Position,
		},
		Status: http.StatusOK,
	})
}
