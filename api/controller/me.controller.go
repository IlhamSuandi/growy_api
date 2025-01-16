package controller

import (
	"net/http"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type meController struct {
	UserUsecase  usecase.UserUsecase
	Logger       *logrus.Logger
	EmailUsecase usecase.EmailUsecase
	MeUsecase    usecase.MeUsecase
}

func NewMeController(
	userUsecase usecase.UserUsecase,
	emailUsecase usecase.EmailUsecase,
	meUsecase usecase.MeUsecase,
) *meController {
	return &meController{
		UserUsecase:  userUsecase,
		Logger:       utils.Log,
		EmailUsecase: emailUsecase,
		MeUsecase:    meUsecase,
	}
}

// @Tags Me
// @Summary get myself informations
// @Description get all informations about myself
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 200 {object} types.Response{data=model.User} "Successfully get myself informations"
// @Router /me [get]
func (uc *meController) GetMe(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(*model.User)

	uc.Logger.Info("[/me] getting user informations")
	user, err := uc.MeUsecase.GetMe(userInfo.Id)
	if err != nil {
		uc.Logger.Errorf("[/me] error getting user informations %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	uc.Logger.Info("[/me] successfully get me")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully get myself informations",
		Data:    user,
		Status:  http.StatusOK,
	})
}

// TODO: make update me
// func (uc *meController) UpdateMe(w http.ResponseWriter, r *http.Request) {
// 	uc.Logger.Info("parsing request body")
// 	var payload dto.UpdateUserRequest
// 	if err := utils.ParseJSON(r, &payload); err != nil {
// 		uc.Logger.Errorf("error parsing request body %s", err)
// 		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
// 			Message: "Error",
// 			Error:   err.Error(),
// 			Status:  http.StatusBadRequest,
// 		})
// 		return
// 	}
//
// 	uc.Logger.Info("getting user informations from auth middleware")
// 	userInfo := r.Context().Value("userInfo").(*model.User)
//
// 	uc.Logger.Info("updating user")
//
// 	updatedUser, err := uc.UserUsecase.UpdateUser(userInfo.Id, model.User{
// 		Username:     *payload.Username,
// 		Email:        *payload.Email,
// 		Password:     *payload.Password,
// 		AuthProvider: *payload.AuthProvider,
// 	})
// 	if err != nil {
// 		uc.Logger.Errorf("error updating user %s", err)
// 		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
// 			Message: "Error",
// 			Status:  http.StatusInternalServerError,
// 		})
// 		return
// 	}
//
// 	response.WriteJSON(w, http.StatusOK, types.Response{
// 		Message: "Successfully update user",
// 		Data: dto.UpdateUserResponse{
// 			Before: dto.UpdateUserRequest{
// 				Username:        &userInfo.Username,
// 				Email:           &userInfo.Email,
// 				IsEmailVerified: &userInfo.IsEmailVerified,
// 				AuthProvider:    &userInfo.AuthProvider,
// 			},
// 			After: dto.UpdateUserRequest{
// 				Username:        &updatedUser.Username,
// 				Email:           &updatedUser.Email,
// 				IsEmailVerified: &updatedUser.IsEmailVerified,
// 				AuthProvider:    &updatedUser.AuthProvider,
// 			},
// 		},
// 		Status: http.StatusOK,
// 	})
// }
