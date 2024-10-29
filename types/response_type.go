package types

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  uint16      `json:"status"`
}

type ErrorResponse struct {
	Error   interface{} `json:"error"`
	Message string      `json:"message" example:"Unauthorized"`
	Status  uint16      `json:"status" example:"401"`
}
