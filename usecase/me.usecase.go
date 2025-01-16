package usecase

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
)

type MeUsecase interface {
	GetMe(userId uint) (*model.User, error)
}

type meUsecase struct {
	meRepo repository.MeRepository
}

func NewMeUsecase(meRepo repository.MeRepository) MeUsecase {
	return &meUsecase{
		meRepo: meRepo,
	}
}

func (mu *meUsecase) GetMe(userId uint) (*model.User, error) {
	return mu.meRepo.GetMe(userId)
}
