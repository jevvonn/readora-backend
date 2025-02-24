package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func New() *log.Logger {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.TraceLevel)

	return logger
}
