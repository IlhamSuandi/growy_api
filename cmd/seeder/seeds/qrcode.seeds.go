package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	qr "github.com/ilhamSuandi/business_assistant/pkg/qrcode"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

func UpdateQrCodes(db *gorm.DB, qrCode *model.QRCode) (*model.QRCode, error) {
	qrRepo := repository.NewQrCodeRepository(db)

	return qrRepo.UpdateQr(qrCode)
}

func GetQrCode(db *gorm.DB, userId uint) (*model.QRCode, error) {
	qrRepo := repository.NewQrCodeRepository(db)

	return qrRepo.GetQrByUserId(userId)
}

func SeedQrCodes(db *gorm.DB) error {
	log := utils.Log
	log.Info("seeding qr codes")
	users := []*model.User{
		&UserOne,
		&UserAdmin,
	}

	for _, user := range users {
		qrData, err := GetQrCode(db, user.Id)
		expiredAt := qr.GetExpirationDate()
		code := qr.GenerateRandomCode(user.UUID)

		qrData.Code = code
		qrData.ExpiresAt = &expiredAt

		if err != nil {
			log.Errorf("error getting qr codes for user %s", user.Email)
			return err
		}

		log.Infof("updating qr codes for user %s", user.Email)
		_, err = UpdateQrCodes(db, qrData)
		if err != nil {
			log.Errorf("error updating qr codes for user %s", user.Email)
			return err
		}
	}

	log.Info("successfully seeded qr codes")
	return nil
}
