package repository

import (
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type BookPostgreSQLItf interface {
	Create(req entity.Book) error
}

type BookPostgreSQL struct {
	db  *gorm.DB
	log logger.LoggerItf
}

func NewBookPostgreSQL(db *gorm.DB, log logger.LoggerItf) BookPostgreSQLItf {
	return &BookPostgreSQL{db, log}
}

func (r *BookPostgreSQL) Create(req entity.Book) error {
	err := r.db.Create(&req).Error
	if err != nil {
		r.log.Error("[BookPostgreSQL][Create]", err)
	}

	return err
}
