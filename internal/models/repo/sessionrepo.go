package repo

import (
	"fmt"

	"github.com/codedbyshoe/goxclone/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

func (s *SessionRepo) CreateSession(session *models.Session) (*models.Session, error) {
	session.SessionID = uuid.New().String()

	result := s.db.Create(session)

	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}

func (s *SessionRepo) GetUserFromSession(sessionID string, userID string) (*models.User, error) {
	var session models.Session

	err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email")
	}).Where("session_id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		return nil, err
	}

	if session.User.ID == 0 {
		return nil, fmt.Errorf("no user associated with the session")
	}

	return &session.User, nil
}
