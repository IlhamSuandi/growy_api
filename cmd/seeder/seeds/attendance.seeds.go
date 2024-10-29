package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

func CheckinUser(db *gorm.DB, userId uint, checkinToken string, location string) error {
	attendanceRepo := repository.NewAttendanceRepository(db)
	qrCodeRepo := repository.NewQrCodeRepository(db)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, qrCodeRepo)

	_, err := attendanceUsecase.CheckInAttendance(checkinToken, userId, location)
	return err
}

func SeedCheckins(db *gorm.DB) error {
	log := utils.Log
	log.Info("seeding checkins")

	users := []model.User{
		UserOne,
		UserAdmin,
	}

	for _, user := range users {
		log.Infof("checking in user %s", user.Email)
		qrData, err := GetQrCode(db, user.Id)
		if err != nil {
			log.Errorf("error getting qr codes for user %s", user.Email)
			return err
		}

		if err := CheckinUser(db, user.Id, qrData.Code, "jakarta"); err != nil {
			log.Errorf("error checking in user %s", user.Email)
			return err
		}
	}

	return nil
}
