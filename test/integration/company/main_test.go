package integration

import (
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
)

var CompanyController *controller.CompanyController

func TestMain(m *testing.M) {
	utils.Log.Info("Starting Company Tests")
  companyRepo := repository.NewCompanyRepository(test.DB)
  companyUsecase := usecase.NewCompanyUsecase(companyRepo)
  CompanyController = controller.NewCompanyController(companyUsecase)

	m.Run()

	utils.Log.Info("Finished Company Tests")
}
