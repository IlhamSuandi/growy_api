package route

import (
	"fmt"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"gorm.io/gorm"
)

func QRCodeRoutes(router *http.ServeMux, db *gorm.DB, path string) {
	qrRepo := repository.NewQrCodeRepository(db)
	userRepo := repository.NewUserRepository(db)
	qrCodeUsecase := usecase.NewQrCodeUsecase(qrRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	controller := controller.NewQrCodeController(qrCodeUsecase, userUsecase)

	router.HandleFunc(fmt.Sprintf("POST %s/{userUUID}", path), controller.CreateQr)
	router.HandleFunc(fmt.Sprintf("GET %s/{userUUID}", path), controller.GetUserQr)
}
