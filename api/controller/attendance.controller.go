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

type AttendanceController struct {
	AttendanceUsecase usecase.AttendanceUsecase
	Logger            *logrus.Logger
}

func NewAttendanceController(attendanceUsecase usecase.AttendanceUsecase) *AttendanceController {
	return &AttendanceController{
		AttendanceUsecase: attendanceUsecase,
		Logger:            utils.Log,
	}
}

// @Tags Attendance
// @Summary check in user
// @Description daily checkin user, use qrcode to check in
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token path string true "token"
// @Param request body dto.CheckInRequest true "request body"
// @Failure 403 {object} types.ErrorResponse "token is empty"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=dto.CheckInResponse} "Successfully Checked In Attendance"
// @Router /attendances/check-in/{token} [post]
func (ac *AttendanceController) CheckIn(w http.ResponseWriter, r *http.Request) {
	var payload dto.CheckInRequest

	ac.Logger.Info("getting checkin token")
	checkinToken := r.PathValue("token")
	if checkinToken == "" {
		ac.Logger.Errorf("token is required")
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "token is required",
			Error:   "token is required",
			Status:  http.StatusForbidden,
		})
		return
	}

	ac.Logger.Info("parsing request body")
	if err := utils.ParseJSON(r, &payload); err != nil {
		ac.Logger.Errorf("error parsing request body %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "error parsing json",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	ac.Logger.Info("getting user info from auth middleware")
	userInfo := r.Context().Value("userInfo").(*model.User)

	ac.Logger.Info("checking in attendance")
	attendance, err := ac.AttendanceUsecase.CheckInAttendance(checkinToken, userInfo.Id, payload.Location)
	if err != nil {
		ac.Logger.Errorf("error checking in attendance %s", err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "error checking in attendance",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	ac.Logger.Info("successfully checked in attendance")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Checked In Attendance",
		Data: dto.CheckInResponse{
			UserId:   userInfo.UUID,
			Status:   attendance.Status,
			CheckIn:  attendance.CheckIn,
			Date:     attendance.Date.Format("2006-01-02"),
			Location: attendance.Location,
		},
		Status: http.StatusOK,
	})
}

// @Tags Attendance
// @Summary Get User Attendances
// @Description get specific user attendances
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token path string true "token"
// @Failure 403 {object} types.ErrorResponse "token is empty"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Success 200 {object} types.Response{data=dto.GetAttendancesResponse} "Successfully got attendances"
// @Router /attendances/{userUUID} [get]
func (ac *AttendanceController) GetUserAttendances(w http.ResponseWriter, r *http.Request) {
	userUUID := r.PathValue("userUUID")

	ac.Logger.Info("parsing user uuid")
	parsedUserId, err := uuid.Parse(userUUID)
	if err != nil {
		ac.Logger.Errorf("error parsing user id %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error parsing UUID",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Info("getting user attendances")
	attendances, err := ac.AttendanceUsecase.GetUserAttendances(parsedUserId)
	if err != nil {
		ac.Logger.Errorf("error getting user attendances %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error getting attendances",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Info("successfully got user attendances")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully got attendances",
		Data: dto.GetAttendancesResponse{
			UserId:      parsedUserId,
			Attendances: attendances,
		},
		Status: http.StatusOK,
	})
}

// check out attendance
func (ac *AttendanceController) CheckOut(w http.ResponseWriter, r *http.Request) {
	ac.Logger.Info("getting user info from auth middleware")
	userInfo := r.Context().Value("userInfo").(*model.User)

	ac.Logger.Info("checking out attendance")
	attendance, err := ac.AttendanceUsecase.CheckOutAttendance(userInfo.Id)
	if err != nil {
		ac.Logger.Errorf("error checking out attendance %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "error checking out attendance",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	ac.Logger.Info("successfully checked out attendance")
	response.WriteJSON(w, http.StatusOK, dto.CheckoutRequest{
		UserId:   userInfo.UUID,
		Status:   attendance.Status,
		CheckIn:  attendance.CheckIn,
		CheckOut: attendance.CheckOut,
		Date:     attendance.Date.Format("2006-01-02"),
		Location: attendance.Location,
	})
}
