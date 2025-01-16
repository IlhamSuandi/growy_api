package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	qr "github.com/ilhamSuandi/business_assistant/pkg/qrcode"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user *model.User) error
	GetUserByUserId(userId uuid.UUID) (*model.User, error)
	GetUsers(page *int, pageSize *int) ([]model.User, error)
	GetUserById(userId uint) (*model.User, error)
	UpdateUser(userId uint, user model.User) (*model.User, error)
	GetUserPermission(userId uint) (*[]model.Permission, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User

	result := ur.db.Preload(clause.Associations).Where("email = ?", email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	tx := ur.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback the transaction in case of panic
		}
	}()

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if a QR code already exists for the user
	var existingQRCode model.QRCode
	if err := tx.Where("id = ?", user.Id).First(&existingQRCode).Error; err == nil {
		// QR code already exists, skip creating a new one
		tx.Commit()
		return nil
	} else if err != gorm.ErrRecordNotFound {
		// An error occurred while querying for the existing QR code
		tx.Rollback()
		return err
	}

	code := qr.GenerateRandomCode(user.UUID)
	expiredAt := qr.GetExpirationDate()

	qrcode := model.QRCode{
		UserId:    user.Id,
		Code:      code,
		ExpiresAt: &expiredAt,
	}

	if err := tx.Save(&qrcode).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ur *userRepository) GetUserByUserId(userId uuid.UUID) (*model.User, error) {
	var user model.User

	result := ur.db.Where("UUID = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, result.Error
	}

	return &user, nil
}

func (ur *userRepository) GetUsers(page *int, pageSize *int) ([]model.User, error) {
	var users []model.User
	defaultPage := 1
	defaultPageSize := 100

	if page == nil || pageSize == nil {
		page = &defaultPage
		pageSize = &defaultPageSize
	}

	offset := (*page - 1) * *pageSize

	result := ur.db.Preload(clause.Associations).Limit(*pageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (ur *userRepository) GetUserById(userId uint) (*model.User, error) {
	var user model.User

	result := ur.db.Preload(clause.Associations).Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *userRepository) UpdateUser(userId uint, updatedUser model.User) (*model.User, error) {
	user, err := ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	user.Email = updatedUser.Email
	user.Username = updatedUser.Username
	user.Password = updatedUser.Password
	user.IsEmailVerified = updatedUser.IsEmailVerified
	user.AuthProvider = updatedUser.AuthProvider

	result := ur.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("error updating user")
	}

	return user, nil
}

func (ur *userRepository) GetUserPermission(userId uint) (*[]model.Permission, error) {
	var permission []model.Permission
	result := ur.db.Where("user_id = ?", userId).Find(&permission)
	if result.Error != nil {
		return nil, result.Error
	}

	return &permission, nil
}
