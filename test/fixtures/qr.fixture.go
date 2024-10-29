package fixtures

import (
	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	qr "github.com/ilhamSuandi/business_assistant/pkg/qrcode"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/test"
)

func CreateQrCode(user model.User) (*model.QRCode, error) {
	qrCodeRepo := repository.NewQrCodeRepository(test.DB)
	code := qr.GenerateRandomCode(user.UUID)
	expiredAt := qr.GetExpirationDate()
	return qrCodeRepo.CreateQr(user.Id, code, expiredAt)
}

func GetQrCode(userId uuid.UUID) (*model.QRCode, error) {
	qrCodeRepo := repository.NewQrCodeRepository(test.DB)
	userRepo := repository.NewUserRepository(test.DB)

	user, err := userRepo.GetUserByUserId(userId)
	if err != nil {
		return nil, err
	}
	return qrCodeRepo.GetQrByUserId(user.Id)
}
