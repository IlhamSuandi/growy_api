package usecase

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type BranchUsecase interface {
	GetCompanyBranches(companyId uint) ([]model.Branch, error)
	GetCompanyBranchByName(companyId uint, branchName string) (*model.Branch, error)
	CreateBranch(branch *model.Branch) error
}

type branchUsecase struct {
	Logger           *logrus.Logger
	BranchRepository repository.BranchRepository
}

func NewBranchUsecase(branchRepository repository.BranchRepository) BranchUsecase {
	return &branchUsecase{
		Logger:           utils.Log,
		BranchRepository: branchRepository,
	}
}

func (eu *branchUsecase) GetCompanyBranches(companyId uint) ([]model.Branch, error) {
	return eu.BranchRepository.GetCompanyBranches(companyId)
}

func (eu *branchUsecase) GetCompanyBranchByName(companyId uint, branchName string) (*model.Branch, error) {
	return eu.BranchRepository.GetCompanyBranchByName(companyId, branchName)
}

func (eu *branchUsecase) CreateBranch(branch *model.Branch) error {
	return eu.BranchRepository.CreateBranch(branch)
}
