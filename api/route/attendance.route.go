package route

import (
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/api/middleware"
	permission "github.com/ilhamSuandi/business_assistant/constant"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"gorm.io/gorm"
)

func AttendanceRoutes(router *http.ServeMux, db *gorm.DB) {
	attendanceRepository := repository.NewAttendanceRepository(db)
	companyRepository := repository.NewCompanyRepository(db)
	qrCodeRepository := repository.NewQrCodeRepository(db)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepository, qrCodeRepository, companyRepository)
	controller := controller.NewAttendanceController(attendanceUsecase)

	router.Handle("GET /attendance", middleware.Permission(
		[]string{
			permission.All,
			permission.Attendances,
		},
		db,
		http.HandlerFunc(controller.GetUserAttendances),
	))

	router.Handle("POST /attendance/check-in", middleware.Permission(
		[]string{
			permission.All,
			permission.Attendances,
		},
		db,
		http.HandlerFunc(controller.CheckIn),
	))

	router.Handle("GET /attendance/check-out", middleware.Permission(
		[]string{
			permission.All,
			permission.Attendances,
		},
		db,
		http.HandlerFunc(controller.CheckOut),
	))
}
