package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"userone@gmail.com"`
	Password string `json:"password" validate:"required" example:"Userone123+"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,exclude_special_char" example:"userone"`
	Password string `json:"password" validate:"required,min=3,contains_uppercase,contains_special_char,contains_num" example:"Userone123+"`
	Email    string `json:"email" validate:"required,email" example:"userone@gmail.com"`
}

type RegisterGoogleUserRequest struct {
	Password string `json:"password" validate:"required,min=3,contains_uppercase,contains_special_char,contains_num" example:"Userone123+"`
}

type LoginGoogleRequest struct {
	Code string `json:"code" validate:"required,min=10"`
}

type RegisteredUserResponse struct {
	Email    string        `json:"email" example:"userone@gmail.com"`
	Username string        `json:"username" example:"userone"`
	Token    TokenResponse `json:"token"`
}

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}
