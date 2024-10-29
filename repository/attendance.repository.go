package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	SaveAttendance(attendance *model.Attendance) (*model.Attendance, error)
	GetUserAttendanceByDate(userId uint, date time.Time) (*model.Attendance, error)
	GetUserAttendances(userId uuid.UUID) ([]*model.Attendance, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{
		db: db,
	}
}

func (ar *attendanceRepository) SaveAttendance(attendance *model.Attendance) (*model.Attendance, error) {
	result := ar.db.Save(&attendance)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("error creating attendance")
	}

	return attendance, nil
}

func (ar *attendanceRepository) GetUserAttendances(userId uuid.UUID) ([]*model.Attendance, error) {
	var user model.User
	var attendances []*model.Attendance

	result := ar.db.Preload("Attendances").Where("UUID = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	attendances = user.Attendances

	return attendances, nil
}

func (ar *attendanceRepository) GetUserAttendanceByDate(userId uint, date time.Time) (*model.Attendance, error) {
	var attendance model.Attendance
	formattedDate := date.Format("2006-01-02")

	result := ar.db.Joins("JOIN user_attendances ua ON ua.attendance_id = attendances.id").
		Where("ua.user_id = ? AND attendances.date = ?", userId, formattedDate).
		Take(&attendance)
	if result.Error != nil {
		return nil, result.Error
	}

	return &attendance, nil
}
