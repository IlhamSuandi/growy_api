package repository

import (
	"errors"
	"time"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type QrCodeRepository interface {
	GetQrByUserId(userId uint) (*model.QRCode, error)
	CreateQr(userId uint, code string, expiredAt time.Time) (*model.QRCode, error)
	UpdateQr(qrcode *model.QRCode) (*model.QRCode, error)
}

type qrcodeRepository struct {
	db *gorm.DB
}

func NewQrCodeRepository(db *gorm.DB) QrCodeRepository {
	return &qrcodeRepository{
		db: db,
	}
}

func (qr *qrcodeRepository) GetQrByUserId(userId uint) (*model.QRCode, error) {
	var qrcode model.QRCode
	result := qr.db.Where("user_id = ?", userId).First(&qrcode)
	if result.Error != nil {
		return nil, result.Error
	}

	return &qrcode, nil
}

func (qr *qrcodeRepository) CreateQr(userId uint, code string, expiredAt time.Time) (*model.QRCode, error) {
	qrcode := model.QRCode{
		UserId:    userId,
		Code:      code,
		ExpiresAt: &expiredAt,
	}

	result := qr.db.Create(&qrcode)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("error creating qr code")
	}

	return &qrcode, nil
}

func (qr *qrcodeRepository) UpdateQr(qrcode *model.QRCode) (*model.QRCode, error) {
	result := qr.db.Save(&qrcode)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("error updating qr code")
	}

	return qrcode, nil
}
