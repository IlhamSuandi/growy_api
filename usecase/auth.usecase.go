package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/auth"
	"github.com/ilhamSuandi/business_assistant/repository"
	types "github.com/ilhamSuandi/business_assistant/types"
)

type AuthUsecase interface {
	GetUserSession(sessionId uuid.UUID) (*model.Session, error)
	CreateSession(session model.Session) error
	IsUserExists(email string) (model.User, bool)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	CreateToken(session types.CreateToken) (string, *types.JwtClaims, error)
	DeleteUserSession(sessionId uuid.UUID) error
}

type authUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) AuthUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (au *authUsecase) GetUserSession(sessionId uuid.UUID) (*model.Session, error) {
	return au.sessionRepo.GetSessionBySessionId(sessionId)
}

func (au *authUsecase) CreateSession(session model.Session) error {
	return au.sessionRepo.SaveSession(session)
}

func (au *authUsecase) IsUserExists(email string) (model.User, bool) {
	user, err := au.userRepo.GetUserByEmail(email)
	return user, err == nil
}

func (au *authUsecase) CreateUser(user *model.User) error {
	// Hash Password
	hashedPass, err := auth.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password %s", err)
	}

	user.Password = string(hashedPass)

	if err := au.userRepo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (au *authUsecase) CreateToken(session types.CreateToken) (string, *types.JwtClaims, error) {
	return auth.CreateToken(session)
}

func (au *authUsecase) DeleteUserSession(sessionId uuid.UUID) error {
	return au.sessionRepo.DeleteSession(sessionId)
}

func (au *authUsecase) UpdateUser(user *model.User) error {
	return au.UpdateUser(user)
}
