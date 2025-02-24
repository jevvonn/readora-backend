package logger

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func New() *log.Logger {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.TraceLevel)

	todayDate := time.Now().Format("2006-01-02")
	path := fmt.Sprintf("logs/%s.log", todayDate)
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger.SetOutput(file)

	return logger
}
