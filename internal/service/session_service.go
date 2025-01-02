package service

import (
	"CloakMail/internal/model"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO: Add caching for sessions to improve performance and reduce database load

const sessionExpiresIn = 30 * 24 * time.Hour

type sessionService struct {
	database *sql.DB
	logger   *logrus.Logger
}

type SessionService interface {
	GenerateRandomToken() (string, error)
	CreateSession(token string, userId int) (model.Session, error)
	ValidateSession(token string) (model.Session, error)
	InvalidateSession(token string) error
}

func NewSessionService(database *sql.DB, logger *logrus.Logger) SessionService {
	return sessionService{
		database: database,
		logger:   logger.WithField("service", "session_service").Logger,
	}
}

func getSessionID(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s sessionService) GenerateRandomToken() (string, error) {
	token := make([]byte, 15)

	_, err := rand.Read(token)
	if err != nil {
		s.logger.Error("Failed to generate token: ", err)
		return "", ErrFailedToGenerateToken
	}

	sessionId := hex.EncodeToString(token)
	return sessionId, nil
}

func (s sessionService) CreateSession(token string, userID int) (model.Session, error) {
	sessionId := getSessionID(token)

	var session model.Session = model.Session{
		ID:        sessionId,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionExpiresIn),
	}

	_, err := s.database.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)", session.ID, session.UserID, session.ExpiresAt)

	if err != nil {
		return model.Session{}, ErrFailedToCreateSession
	}

	return session, nil
}

func (s sessionService) ValidateSession(token string) (model.Session, error) {
	sessionId := getSessionID(token)

	var session model.Session

	row := s.database.QueryRow("SELECT id, user_id, expires_at FROM sessions WHERE id = $1", sessionId)

	if err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt); err != nil {
		s.logger.Error("Failed to validate session (scanning row): ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return model.Session{}, ErrSessionDoesNotExist
		}
		return model.Session{}, ErrFailedGetSession
	}

	now := time.Now()

	if now.After(session.ExpiresAt) {
		_, err := s.database.Exec("DELETE FROM sessions WHERE id = $1", sessionId)
		if err != nil {
			s.logger.Error("Failed to delete session: ", err)
		}
		return model.Session{}, ErrSessionExpired
	}

	if now.After(session.ExpiresAt.Add(-15 * 24 * time.Hour)) {
		session.ExpiresAt = now.Add(sessionExpiresIn)
		_, err := s.database.Exec("UPDATE sessions SET expires_at = $1 WHERE id = $2", session.ExpiresAt, sessionId)
		if err != nil {
			s.logger.Error("Failed to update session: ", err)
		}
	}

	return session, nil

}

func (s sessionService) InvalidateSession(sessionID string) error {
	_, err := s.database.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		s.logger.Error("Failed to invalidate session: ", err)
		return ErrFailedInvalidateSession
	}
	return nil
}
