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

func BranchRoutes(
	router *http.ServeMux,
	db *gorm.DB,
) {
	branchRepo := repository.NewBranchRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	branchUsecase := usecase.NewBranchUsecase(branchRepo)
	companyUsecase := usecase.NewCompanyUsecase(companyRepo)
	controller := controller.NewBranchController(branchUsecase, companyUsecase)

	permissions := []string{
		permission.All,
	}

	router.Handle("GET /branch", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.GetCompanyBranches),
	))

	router.Handle("POST /branch", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(controller.CreateBranch),
	))
}
