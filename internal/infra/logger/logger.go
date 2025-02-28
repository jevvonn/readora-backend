package logger

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type LoggerItf interface {
	Info(path string, message string)
	Error(path string, err error)
	Warn(path string, message string)
}

type Logger struct {
	log *log.Logger
}

func New() LoggerItf {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.TraceLevel)

	todayDate := time.Now().Format("2006-01-02")
	path := fmt.Sprintf("logs/%s.log", todayDate)
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger.SetOutput(file)

	return &Logger{logger}
}

func (l *Logger) Info(path string, message string) {
	l.log.WithFields(log.Fields{
		"path": path,
	}).Info(message)
}

func (l *Logger) Error(path string, err error) {
	l.log.WithFields(log.Fields{
		"path": path,
	}).Error(err)
}

func (l *Logger) Warn(path string, message string) {
	l.log.WithFields(log.Fields{
		"path": path,
	}).Warn(message)
}
