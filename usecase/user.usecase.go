package usecase

import (
	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
)

type UserUsecase interface {
	GetUserByEmail(email string) (model.User, error)
	GetUserByUserId(userId uuid.UUID) (*model.User, error)
	GetUsers() ([]model.User, error)
	CreateUser(user *model.User) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) CreateUser(user *model.User) error {
	return uu.userRepo.CreateUser(user)
}

func (uu *userUsecase) GetUserByUserId(userId uuid.UUID) (*model.User, error) {
	return uu.userRepo.GetUserByUserId(userId)
}

// find by email and return user or error
func (uu *userUsecase) GetUserByEmail(email string) (model.User, error) {
	return uu.GetUserByEmail(email)
}

func (uu *userUsecase) GetUsers() ([]model.User, error) {
	return uu.userRepo.GetUsers(nil, nil)
}
