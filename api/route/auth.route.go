package route

import (
	"fmt"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"gorm.io/gorm"
)

func AuthRoutes(router *http.ServeMux, db *gorm.DB, path string) {
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	useCase := usecase.NewAuthUsecase(userRepository, sessionRepository)

	controller := controller.NewAuthController(useCase)

	router.HandleFunc(fmt.Sprintf("POST %s/register", path), controller.Register)
	router.HandleFunc(fmt.Sprintf("POST %s/login", path), controller.Login)
	router.HandleFunc(fmt.Sprintf("POST %s/token/renew", path), controller.RenewAccessToken)
	router.HandleFunc(fmt.Sprintf("POST %s/logout", path), controller.Logout)
	router.HandleFunc(fmt.Sprintf("GET %s/google/login", path), controller.GoogleLogin)
	router.HandleFunc(fmt.Sprintf("GET %s/google/callback", path), controller.GoogleCallback)
}
