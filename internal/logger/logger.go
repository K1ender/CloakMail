package logger

import (
	"CloakMail/internal/config"

	log "github.com/sirupsen/logrus"
)

func Init(cfg *config.Config) *log.Logger {
	var logger *log.Logger

	if cfg.Env == config.ProdEnv {
		logger = log.New()
		logger.Level = log.InfoLevel
		logger.Formatter = &log.JSONFormatter{}
	} else {
		logger = log.New()
		logger.Level = log.DebugLevel
		logger.Formatter = &log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		}
	}
	return logger
}
