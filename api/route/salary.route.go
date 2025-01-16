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

func SalaryRoutes(
	router *http.ServeMux,
	db *gorm.DB,
) {
	// TODO: implement salary routes

	repo := repository.NewSalaryRepository(db)
	usecase := usecase.NewSalaryUsecase(repo)
	controller := controller.NewSalaryController(usecase)

	permissions := []string{
		permission.Salary,
	}

	router.Handle("POST /salary", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.CreateEmployeeSalary),
	))
}
