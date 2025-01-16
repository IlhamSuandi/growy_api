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

func EmployeeRoutes(router *http.ServeMux, db *gorm.DB) {
	employeeRepo := repository.NewEmployeeRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	userRepo := repository.NewUserRepository(db)
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo)
	companyUsecase := usecase.NewCompanyUsecase(companyRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	branchrepo := repository.NewBranchRepository(db)
	branchUsecase := usecase.NewBranchUsecase(branchrepo)
	employeeController := controller.NewEmployeeController(
		employeeUsecase,
		userUsecase,
		companyUsecase,
		branchUsecase,
	)

	permissions := []string{
		permission.EmployeeRoute,
	}

	// Add Employee
	router.Handle("POST /employee", middleware.Permission(
		permissions,
		db,
		http.HandlerFunc(employeeController.AddEmployee),
	))
}
