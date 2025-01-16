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

func UserRoutes(router *http.ServeMux, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	usecase := usecase.NewUserUsecase(userRepository)
	controller := controller.NewUserController(usecase)
	permissions := []string{
		permission.User,
	}

	router.Handle("GET /users", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.GetUsers),
	))

	router.Handle("GET /users/{userUUID}", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.GetUserId),
	))
}
