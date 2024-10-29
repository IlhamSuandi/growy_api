package route

import (
	"fmt"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"gorm.io/gorm"
)

func UserRoutes(router *http.ServeMux, db *gorm.DB, path string) {
	userRepository := repository.NewUserRepository(db)
	usecase := usecase.NewUserUsecase(userRepository)
	controller := controller.NewUserController(usecase)

	router.HandleFunc(fmt.Sprintf("GET %s", path), controller.GetUsers)
	router.HandleFunc(fmt.Sprintf("GET %s/{userID}", path), controller.GetUserId)
}
