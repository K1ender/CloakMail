package service

import (
	"CloakMail/internal/model"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	database *sql.DB
	logger   *logrus.Logger
}

type UserService interface {
	CreateUser(email string, password string) (int, error)
	GetUserById(id int) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	ChangePassword(id int, newPassword string) error
}

func NewUserService(database *sql.DB, logger *logrus.Logger) UserService {
	return userService{
		database: database,
		logger:   logger.WithField("service", "user_service").Logger,
	}
}

func (s userService) CreateUser(email string, password string) (int, error) {
	s.logger.WithFields(logrus.Fields{"email": email}).Debug("Creating user with email")

	tx, err := s.database.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction: ", err)
		return 0, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, ErrFailedToHashPassword
	}

	row := tx.QueryRow("INSERT INTO users(email, password) values ($1, $2) RETURNING id", email, string(hashedPassword))

	var userID int
	if err := row.Scan(&userID); err != nil {
		s.logger.Error("Failed to create user (scanning row): ", err)
		return 0, ErrFailedToCreateUser
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		return 0, ErrFailedToCommitTransaction
	}
	return userID, nil
}

func (s userService) ChangePassword(id int, newPassword string) error {
	s.logger.WithFields(logrus.Fields{"id": id}).Debug("Changing password")

	tx, err := s.database.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction: ", err)
		return ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return ErrFailedToHashPassword
	}

	_, err = tx.Exec("UPDATE users SET password=$1 WHERE id=$2", hashedPassword, id)
	if err != nil {
		s.logger.Error("Failed to change password: ", err)
		return ErrFailedChangePassword
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		return ErrFailedToCommitTransaction
	}

	return nil
}

func (s userService) GetUserById(id int) (model.User, error) {
	user := model.User{}

	row := s.database.QueryRow("SELECT id, email, password FROM users WHERE id = $1", id)

	if err := row.Scan(&user.ID, &user.Email, &user.HashedPassword); err != nil {
		s.logger.Error("Failed to get user (scanning row): ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, ErrFailedToFindUser
	}

	return user, nil
}

func (s userService) GetUserByEmail(email string) (model.User, error) {
	user := model.User{}

	row := s.database.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email)

	if err := row.Scan(&user.ID, &user.Email, &user.HashedPassword); err != nil {
		s.logger.Error("Failed to get user (scanning row): ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, ErrFailedToFindUser
	}

	return user, nil
}
