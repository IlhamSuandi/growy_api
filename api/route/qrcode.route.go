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

func QRCodeRoutes(
	router *http.ServeMux,
	db *gorm.DB,
) {
	qrRepo := repository.NewQrCodeRepository(db)
	userRepo := repository.NewUserRepository(db)
	qrCodeUsecase := usecase.NewQrCodeUsecase(qrRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	controller := controller.NewQrCodeController(qrCodeUsecase, userUsecase)

	router.Handle("GET /qrcode", middleware.Permission(
		[]string{
			permission.All,
			permission.Qrcode,
		},
		db,
		http.HandlerFunc(controller.GetUserQr),
	))

	// router.Handle("GET /qrcode/{userId}", middleware.Permission(
	// 	[]string{
	// 		permission.All,
	// 		permission.Qrcode,
	// 	},
	// 	db,
	// 	http.HandlerFunc(controller.GetUserQr),
	// ))

	router.Handle("POST /qrcode", middleware.Permission(
		[]string{
			permission.Owner,
			permission.Qrcode,
		},
		db,
		http.HandlerFunc(controller.CreateQr),
	))
}
