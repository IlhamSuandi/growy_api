package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceUsecase interface {
	CheckInAttendance(checkinToken string, userId uint, location string) (*model.Attendance, error)
	GetUserAttendances(userId uuid.UUID) ([]*model.Attendance, error)
	CheckOutAttendance(userId uint) (*model.Attendance, error)
}

type attendanceUsecase struct {
	logger         *logrus.Logger
	attendanceRepo repository.AttendanceRepository
	qrCodeRepo     repository.QrCodeRepository
}

func NewAttendanceUsecase(
	attendanceRepo repository.AttendanceRepository,
	qrCodeRepo repository.QrCodeRepository,
) AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: attendanceRepo,
		qrCodeRepo:     qrCodeRepo,
		logger:         utils.Log,
	}
}

func (au *attendanceUsecase) CheckInAttendance(
	checkinToken string,
	userId uint,
	location string,
) (*model.Attendance, error) {
	// get current time
	au.logger.Info("getting current time")
	now := time.Now()

	// check check-in token
	qrData, err := au.qrCodeRepo.GetQrByUserId(userId)
	if err != nil {
		return nil, err
	}

	// check if qr code is used
	if qrData.IsUsed {
		return nil, errors.New("you have already checked in")
	}

	// check if qr is valid
	if qrData.Code != checkinToken {
		return nil, errors.New("checkin token is invalid")
	}

	// check if qr is expired
	if qrData.ExpiresAt.Before(now) {
		return nil, errors.New("qr code is expired")
	}

	// check if attendance already exists
	existingAttendance, err := au.attendanceRepo.GetUserAttendanceByDate(userId, now)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		au.logger.Error(err)
		return nil, err
	}

	if existingAttendance != nil {
		au.logger.Error("Attendance already exists")
		return nil, errors.New("attendance already exists")
	}

	// TODO: check user checkin time and add status to it
	status := model.AttendanceStatusPresent

	// creating checkin newAttendance
	attendance, err := au.attendanceRepo.SaveAttendance(&model.Attendance{
		Users:    []*model.User{{Model: model.Model{Id: userId}}},
		CheckIn:  &now,
		Date:     now,
		Location: &location,
		Status:   status,
	})
	if err != nil {
		return nil, err
	}

	// update qr code
	qrData.IsUsed = true

	// update qr code
	qrData, err = au.qrCodeRepo.UpdateQr(qrData)
	if err != nil {
		return nil, err
	}

	au.logger.Info("Checkin attendance created successfully")
	return attendance, nil
}

func (au *attendanceUsecase) GetUserAttendances(userId uuid.UUID) ([]*model.Attendance, error) {
	return au.attendanceRepo.GetUserAttendances(userId)
}

func (au *attendanceUsecase) CheckOutAttendance(userId uint) (*model.Attendance, error) {
	// get time to check attendance by date
	now := time.Now()
	// check if attendance is exist
	existingAttendance, err := au.attendanceRepo.GetUserAttendanceByDate(userId, now)
	if err != nil {
		return nil, errors.New("attendance does not exist")
	}

	if existingAttendance.CheckOut != nil {
		return nil, errors.New("user already checkout")
	}

	// add checkout to current time
	existingAttendance.CheckOut = &now

	// update attendance with checkout
	updatedAttendance, err := au.attendanceRepo.SaveAttendance(existingAttendance)
	if err != nil {
		return nil, fmt.Errorf("error updating attendance %s", err)
	}

	return updatedAttendance, nil
}
