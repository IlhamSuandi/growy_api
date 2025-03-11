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

func CompanyRoutes(router *http.ServeMux, db *gorm.DB) {
	companyRepo := repository.NewCompanyRepository(db)
	companyUsecase := usecase.NewCompanyUsecase(companyRepo)
	controller := controller.NewCompanyController(companyUsecase)

	permissions := []string{
		permission.All,
	}

	// Create Company
	router.Handle("POST /company", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.CreateCompany),
	))

	// get user companies
	router.Handle("GET /company", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.GetUserCompanies),
	))
}
