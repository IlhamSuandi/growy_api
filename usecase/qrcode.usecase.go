package usecase

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	qr "github.com/ilhamSuandi/business_assistant/pkg/qrcode"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QrCodeUsecase interface {
	CreateQrCode(user model.User) ([]byte, *model.QRCode, error)
	GetUserQr(userId uint) ([]byte, *model.QRCode, error)
}

type qrcodeUsecase struct {
	qrcodeRepo repository.QrCodeRepository
	userRepo   repository.UserRepository
	logger     *logrus.Logger
}

func NewQrCodeUsecase(qrcodeRepo repository.QrCodeRepository, userRepo repository.UserRepository) QrCodeUsecase {
	return &qrcodeUsecase{
		qrcodeRepo: qrcodeRepo,
		userRepo:   userRepo,
		logger:     utils.Log,
	}
}

func (qu *qrcodeUsecase) CreateQrCode(user model.User) ([]byte, *model.QRCode, error) {
	// find qr code by user id
	qu.logger.Info("generating random code")
	code := qr.GenerateRandomCode(user.UUID)
	expiredAt := qr.GetExpirationDate()

	qu.logger.Info("generating qrcode from random code")
	qrCode, err := qr.GenerateQrCode(code)
	if err != nil {
		qu.logger.Errorf("error generating qrcode from random code %s", err)
		return nil, nil, err
	}

	qu.logger.Info("getting qrcode by user id")
	existedQrCode, err := qu.qrcodeRepo.GetQrByUserId(user.Id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		qu.logger.Errorf("error getting qrcode by user id %s", err)
		return nil, nil, err
	}

	if existedQrCode != nil {
		qu.logger.Info("updating qrcode")
		existedQrCode.Code = code
		existedQrCode.ExpiresAt = &expiredAt
		existedQrCode.IsUsed = false

		qrData, err := qu.qrcodeRepo.UpdateQr(existedQrCode)
		return qrCode, qrData, err
	}

	qu.logger.Info("creating qrcode")
	qrData, err := qu.qrcodeRepo.CreateQr(user.Id, code, expiredAt)
	return qrCode, qrData, err
}

func (qu *qrcodeUsecase) GetUserQr(userId uint) ([]byte, *model.QRCode, error) {
	qu.logger.Info("getting qrcode by user id")
	qrData, err := qu.qrcodeRepo.GetQrByUserId(userId)
	if err != nil {
		qu.logger.Errorf("error getting qrcode by user id %s", err)
		return nil, nil, err
	}

	qu.logger.Info("generating qrcode from qrcode code")
	qrCode, err := qr.GenerateQrCode(qrData.Code)
	if err != nil {
		qu.logger.Errorf("error generating qrcode from qrcode code %s", err)
		return nil, nil, err
	}

	return qrCode, qrData, nil
}
