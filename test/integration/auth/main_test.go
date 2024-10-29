package integration

import (
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/controller"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
)

var authController *controller.AuthController

func TestMain(m *testing.M) {
	utils.Log.Info("Starting Auth Tests")

	userRepo := repository.NewUserRepository(test.DB)
	sessionRepo := repository.NewSessionRepository(test.DB)
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo)
	authController = controller.NewAuthController(authUsecase)

	m.Run()

	utils.Log.Info("Finished Auth Tests")
}
