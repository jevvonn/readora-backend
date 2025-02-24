package repository

import (
	"github.com/jevvonn/reodora-backend/internal/domain/entity"
	"github.com/jevvonn/reodora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type UserPostgreSQLItf interface {
	GetSpecificUser(user entity.User) (entity.User, error)
	CreateUser(user entity.User) error
}

type UserPostgreSQL struct {
	db  *gorm.DB
	log logger.LoggerItf
}

func NewUserPostgreSQL(db *gorm.DB, log logger.LoggerItf) UserPostgreSQLItf {
	return &UserPostgreSQL{db, log}
}

func (r *UserPostgreSQL) GetSpecificUser(user entity.User) (entity.User, error) {
	log := "[UserPostgreSQL][GetSpecificUser]"

	var result entity.User
	err := r.db.First(&result, &user).Error

	if err != nil {
		r.log.Error(log, err)
	}

	return result, err
}

func (r *UserPostgreSQL) CreateUser(user entity.User) error {
	return r.db.Create(&user).Error
}
