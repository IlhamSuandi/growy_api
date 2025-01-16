package dto

import "github.com/ilhamSuandi/business_assistant/database/model"

type CreateCompanyRequest struct {
	Name    string `json:"name" example:"growy" validate:"min=3"`
	Address string `json:"address" example:"jakarta" validate:"min=5"`
}

type CreateCompanyResponse struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	OwnerEmail string `json:"owner_email"`
}

type GetUserCompaniesResponse struct {
	Company []model.Company `json:"company"`
}
