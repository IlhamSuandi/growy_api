package integration

import (
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
)

var attendanceController *controller.AttendanceController

func TestMain(m *testing.M) {
	utils.Log.Info("Starting Auth Tests")

	attendanceRepo := repository.NewAttendanceRepository(test.DB)
  companyRepo := repository.NewCompanyRepository(test.DB)
	qrCodeRepo := repository.NewQrCodeRepository(test.DB)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, qrCodeRepo, companyRepo)
	attendanceController = controller.NewAttendanceController(attendanceUsecase)

	m.Run()

	utils.Log.Info("Finished Auth Tests")
}
