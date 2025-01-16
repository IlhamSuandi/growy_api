package dto

import "github.com/google/uuid"

type CreateQrRequest struct {
	UserId uuid.UUID `validate:"required" json:"user_id"`
}

type CreateQrResponse struct {
	Id        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	ExpiresAt string    `json:"expires_at"`
	QrCode    []byte    `json:"qr_code"`
}

type GetQrResponse struct {
	Id        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	ExpiresAt string    `json:"expires_at"`
	QrCode    []byte    `json:"qr_code"`
}
