package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type UserPostgreSQLItf interface {
	GetSpecificUser(user entity.User) (entity.User, error)
	GetUserByEmailOrUsername(email string, username string) (entity.User, error)
	CreateUser(user entity.User) error
	UpdateUser(user entity.User) error
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

func (r *UserPostgreSQL) GetUserByEmailOrUsername(email string, username string) (entity.User, error) {
	log := "[UserPostgreSQL][GetUserByEmailOrUsername]"

	var result entity.User
	err := r.db.Where("email = ? OR username = ?", email, username).First(&result).Error

	if err != nil {
		r.log.Error(log, err)
	}

	return result, err
}

func (r *UserPostgreSQL) CreateUser(user entity.User) error {
	return r.db.Create(&user).Error
}

func (r *UserPostgreSQL) UpdateUser(user entity.User) error {
	log := "[UserPostgreSQL][UpdateUser]"

	if user.ID == uuid.Nil {
		return errors.New("user ID is required")
	}

	err := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(&user).Error

	if err != nil {
		r.log.Error(log, err)
	}

	return err
}
