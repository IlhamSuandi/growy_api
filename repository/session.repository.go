package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type SessionRepository interface {
	GetSessionBySessionId(sessionId uuid.UUID) (*model.Session, error)
	SaveSession(session model.Session) error
	DeleteSession(sessionId uuid.UUID) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (sr *sessionRepository) GetSessionBySessionId(sessionId uuid.UUID) (*model.Session, error) {
	var session model.Session

	result := sr.db.Where("uuid = ?", sessionId).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

func (sr *sessionRepository) SaveSession(session model.Session) error {
	if err := sr.db.Save(&session).Error; err != nil {
		return err
	}

	return nil
}

func (sr *sessionRepository) DeleteSession(sessionId uuid.UUID) error {
	result := sr.db.Where("session_id = ?", sessionId).Delete(&model.Session{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Failed to delete session")
	}

	return nil
}
