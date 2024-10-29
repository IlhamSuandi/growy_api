package route

import (
	"fmt"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"gorm.io/gorm"
)

func CompanyRoutes(router *http.ServeMux, db *gorm.DB, path string) {
	attendanceRepository := repository.NewAttendanceRepository(db)
	qrCodeRepository := repository.NewQrCodeRepository(db)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepository, qrCodeRepository)
	controller := controller.NewAttendanceController(attendanceUsecase)

	router.HandleFunc(fmt.Sprintf("POST %s", path), controller.GetUserAttendances)
	router.HandleFunc(fmt.Sprintf("POST %s/check-in/{token}", path), controller.CheckIn)
	router.HandleFunc(fmt.Sprintf("PATCH %s/check-out", path), controller.CheckOut)
}
