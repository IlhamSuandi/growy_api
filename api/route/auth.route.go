package route

import (
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"gorm.io/gorm"
)

func AuthRoutes(router *http.ServeMux, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	useCase := usecase.NewAuthUsecase(userRepository, sessionRepository)
	controller := controller.NewAuthController(useCase)

	router.HandleFunc("GET /auth/token/renew", controller.RenewAccessToken)
	router.HandleFunc("GET /auth/google/login", controller.GoogleLogin)
	router.HandleFunc("GET /auth/google/callback", controller.GoogleCallback)
	router.HandleFunc("POST /auth/register", controller.Register)
	router.HandleFunc("POST /auth/login", controller.Login)
	router.HandleFunc("POST /auth/logout", controller.Logout)
}
