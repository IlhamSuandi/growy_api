package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

func CheckinUser(db *gorm.DB, userId uint, location string) error {
	attendanceRepo := repository.NewAttendanceRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	qrCodeRepo := repository.NewQrCodeRepository(db)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, qrCodeRepo, companyRepo)

	_, err := attendanceUsecase.CheckInAttendance(userId, location)
	return err
}

func SeedCheckins(db *gorm.DB) {
	log := utils.Log
	log.Info("seeding checkins")

	users := []model.User{
		UserOne,
		UserAdmin,
	}

	for _, user := range users {
		log.Infof("checking in user %s", user.Email)

		if err := CheckinUser(db, user.Id, "jakarta"); err != nil {
			log.Errorf("error checking in user %s", user.Email)
			log.Fatal(err)
		}
	}

}
