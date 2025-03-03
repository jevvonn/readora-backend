package repository

import (
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type BookPostgreSQLItf interface{}

type BookPostgreSQL struct {
	db  *gorm.DB
	log logger.LoggerItf
}

func NewBookPostgreSQL(db *gorm.DB, log logger.LoggerItf) BookPostgreSQLItf {
	return &BookPostgreSQL{db, log}
}
