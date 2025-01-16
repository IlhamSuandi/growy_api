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
// @Param request body dto.CreateQrRequest true "Request body"
// @Failure 500 {object} types.ErrorResponse "error creating qrcode"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=dto.CreateQrResponse} "successfully created qrcode"
// @Router /qrcode [post]
func (qc *QrCodeController) CreateQr(w http.ResponseWriter, r *http.Request) {
	qc.Logger.Infof("[%s /qrcode] parsing request body", r.Method)

	var payload dto.CreateQrRequest
	// get the request body
	if err := utils.ParseJSON(r, &payload); err != nil {
		qc.Logger.Errorf("[%s /qrcode] error parsing request body %s", r.Method, err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error Parsing Request Body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	qc.Logger.Infof("[%s /qrcode] getting user by user id", r.Method)
	user, err := qc.UserUsecase.GetUserByUserId(payload.UserId)
	if err != nil {
		qc.Logger.Errorf("[%s /qrcode] error getting user by user id %s", r.Method, err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "User not found",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	qc.Logger.Infof("[%s /qrcode] creating user qrcode", r.Method)
	qrCode, qrData, err := qc.QrCodeUsecase.CreateQrCode(*user)
	if err != nil {
		qc.Logger.Errorf("[%s /qrcode] error creating user qrcode %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error creating qrcode",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	qc.Logger.Infof("[%s /qrcode] successfully created user qrcode", r.Method)
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
// @Param userId query string true "userId"
// @Failure 500 {object} types.ErrorResponse "error getting qrcode"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=dto.CreateQrResponse} "successfully created qrcode"
// @Router /qrcode [get]
func (qc *QrCodeController) GetUserQr(w http.ResponseWriter, r *http.Request) {
	qc.Logger.Infof("[%s /qrcode] getting user uuid from path", r.Method)
	// pathUserId := uuid.MustParse(r.PathValue("userId"))
	pathUserId := r.URL.Query().Get("userId")
	if pathUserId == "" {
		qc.Logger.Errorf("[%s /qrcode] token is required", r.Method)
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "token is required",
			Error:   "token is required",
			Status:  http.StatusForbidden,
		})
		return
	}

	qc.Logger.Infof("[%s /qrcode] parsing user uuid", r.Method)
	parsedUserId, err := uuid.Parse(pathUserId)

	qc.Logger.Infof("[%s /qrcode] getting user info from auth middleware", r.Method)
	userInfo := r.Context().Value("userInfo").(*model.User)

	if parsedUserId != userInfo.UUID {
		qc.Logger.Errorf("[%s %s /qrcode] user uuid is not equal", r.Method, r.Method)
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "User UUID is not equal",
			Error:   "User UUID is not equal",
			Status:  http.StatusForbidden,
		})
		return
	}

	qc.Logger.Infof("[%s /qrcode] getting user qrcode", r.Method)
	qrCode, qrData, err := qc.QrCodeUsecase.GetUserQr(userInfo.Id)
	if err != nil {
		qc.Logger.Errorf("[%s /qrcode] error getting user qrcode %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error getting qrcode",
			Error:   err,
			Status:  http.StatusInternalServerError,
		})
		return
	}

	qc.Logger.Infof("[%s /qrcode] successfully get user qrcode", r.Method)
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
