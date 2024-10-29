package qr

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/skip2/go-qrcode"
)

func GetExpirationDate() time.Time {
	return time.Now().Add(time.Hour * 12)
}

func GenerateRandomCode(userId uuid.UUID) string {
	return uuid.New().String() + userId.String()
}

func GenerateQrCode(code string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", config.URL, code)
	png, err := qrcode.Encode(url, qrcode.High, 512)
	if err != nil {
		return nil, err
	}

	return png, nil
}
