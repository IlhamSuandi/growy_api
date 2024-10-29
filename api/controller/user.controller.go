package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	types "github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type userController struct {
	UserUsecase usecase.UserUsecase
	Logger      *logrus.Logger
}

func NewUserController(userUsecase usecase.UserUsecase) *userController {
	return &userController{
		UserUsecase: userUsecase,
		Logger:      utils.Log,
	}
}

// @Tags User
// @Summary Get User by UserId
// @Description getting specific user information
// @Param userID path string true "User ID"
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Failure 400 {object} types.ErrorResponse "Bad Request"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=model.User} "Successfully get all users"
// @Router /users/{userID} [get]
func (uc *userController) GetUserId(w http.ResponseWriter, r *http.Request) {
	uc.Logger.Info("parsing user id")
	userId, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		uc.Logger.Errorf("error parsing user id %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	uc.Logger.Info("getting user by user id")
	user, err := uc.UserUsecase.GetUserByUserId(userId)
	if err != nil {
		uc.Logger.Errorf("error getting user by user id %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	uc.Logger.Info("successfully get user by user id")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Get User",
		Data:    user,
		Status:  http.StatusOK,
	})
}

// @Tags User
// @Summary Get All Users
// @Description getting all users information
// @Router /users [get]
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=model.User} "Successfully get all users"
func (h *userController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserUsecase.GetUsers()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully get all users",
		Data:    users,
		Status:  http.StatusOK,
	})
}
