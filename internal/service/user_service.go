package service

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type userService struct {
	database *sql.DB
	logger   *logrus.Logger
}

type UserService interface{}

func NewUserService(database *sql.DB, logger *logrus.Logger) UserService {
	return userService{
		database: database,
		logger:   logger.WithField("service", "user_service.go").Logger,
	}
}
