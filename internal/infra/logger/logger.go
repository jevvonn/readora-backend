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

	// Check path
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Check log file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(path)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

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
