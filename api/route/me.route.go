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

func MeRoutes(
	router *http.ServeMux,
	db *gorm.DB,
) {
	userRepository := repository.NewUserRepository(db)
	meRepository := repository.NewMeRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	emailUsecase := usecase.NewEmailUsecase()
	meUsecase := usecase.NewMeUsecase(meRepository)
	controller := controller.NewMeController(userUsecase, emailUsecase, meUsecase)

	permissions := []string{
		permission.All,
	}

	// TODO: get me
	router.Handle("GET /me", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.GetMe),
	))

	// TODO: update me
	// router.Handle("PATCH /me", middleware.Permission(
	// 	permissions,
	// 	db,
	// 	http.HandlerFunc(controller.UpdateMe),
	// ))
}
