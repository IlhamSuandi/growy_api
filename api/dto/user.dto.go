package dto

type UpdateUserRequest struct {
	Email           *string `json:"email" validate:"omitempty,email" example:"johndoe@gmail.com"`
	Username        *string `json:"username" validate:"omitempty,min=3,exclude_special_char" example:"john doe"`
	Password        *string `json:"password" validate:"omitempty,min=3,contains_uppercase,contains_special_char,contains_num" example:"John123+"`
	IsEmailVerified *bool   `json:"verify_email" validate:"omitempty,boolean" example:"true"`
	AuthProvider    *string `json:"auth_provider" validate:"-" example:"google"`
}

type UpdateUserResponse struct {
	Before UpdateUserRequest `json:"before"`
	After  UpdateUserRequest `json:"after"`
}
