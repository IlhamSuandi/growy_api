package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type QrCodeController struct {
	QrCodeUsecase usecase.QrCodeUsecase
	UserUsecase   usecase.UserUsecase
	Logger        *logrus.Logger
}

func NewQrCodeController(qrcodeUsecase usecase.QrCodeUsecase, userUsecase usecase.UserUsecase) *QrCodeController {
	return &QrCodeController{
		QrCodeUsecase: qrcodeUsecase,
		UserUsecase:   userUsecase,
		Logger:        utils.Log,
	}
}

// @Tags QrCode
// @Summary re-generate user qrcode
// @Description manually re-generating user qrcode
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "userId"
// @Failure 500 {object} types.ErrorResponse "error creating qrcode"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=dto.CreateQrResponse} "successfully created qrcode"
// @Router /qrcode/{userId} [post]
func (qc *QrCodeController) CreateQr(w http.ResponseWriter, r *http.Request) {
	qc.Logger.Info("getting user uuid from path")
	userId := uuid.MustParse(r.PathValue("userUUID"))

	qc.Logger.Info("getting user by user id")
	user, err := qc.UserUsecase.GetUserByUserId(userId)
	if err != nil {
		qc.Logger.Errorf("error getting user by user id %s", err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "User not found",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	qc.Logger.Info("creating user qrcode")
	qrCode, qrData, err := qc.QrCodeUsecase.CreateQrCode(*user)
	if err != nil {
		qc.Logger.Errorf("error creating user qrcode %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error creating qrcode",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	qc.Logger.Info("successfully created user qrcode")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "successfully created qrcode",
		Data: dto.CreateQrResponse{
			Id:        qrData.UUID,
			Code:      qrData.Code,
			ExpiresAt: qrData.ExpiresAt.String(),
			QrCode:    qrCode,
		},
		Status: http.StatusOK,
	})
}

// @Tags QrCode
// @Summary check in user
// @Description daily checkin user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "userId"
// @Failure 500 {object} types.ErrorResponse "error getting qrcode"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=dto.CreateQrResponse} "successfully created qrcode"
// @Router /qrcode/{userId} [get]
func (qc *QrCodeController) GetUserQr(w http.ResponseWriter, r *http.Request) {
	qc.Logger.Info("getting user uuid from path")
	pathUserId := uuid.MustParse(r.PathValue("userUUID"))
	userInfo := r.Context().Value("userInfo").(*model.User)

	if pathUserId != userInfo.UUID {
		qc.Logger.Errorf("user uuid is not equal")
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "User UUID is not equal",
			Error:   "User UUID is not equal",
			Status:  http.StatusForbidden,
		})
	}

	qc.Logger.Info("getting user qrcode")
	qrCode, qrData, err := qc.QrCodeUsecase.GetUserQr(userInfo.Id)
	if err != nil {
		qc.Logger.Errorf("error getting user qrcode %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error getting qrcode",
			Error:   err,
			Status:  http.StatusInternalServerError,
		})
		return
	}

	qc.Logger.Info("successfully get user qrcode")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "successfully get qrcode",
		Data: dto.GetQrResponse{
			Id:        qrData.UUID,
			Code:      qrData.Code,
			ExpiresAt: qrData.ExpiresAt.String(),
			QrCode:    qrCode,
		},
		Status: http.StatusOK,
	})
}
