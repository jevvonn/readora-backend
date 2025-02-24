package repository

import (
	"github.com/jevvonn/reodora-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type UserPostgreSQLItf interface{}

type UserPostgreSQL struct {
	db *gorm.DB
}

func NewUserPostgreSQL(db *gorm.DB) UserPostgreSQLItf {
	return &UserPostgreSQL{db}
}

func (r *UserPostgreSQL) GetSpecificUser(user entity.User) error {
	return r.db.First(&user).Error
}

func (r *UserPostgreSQL) CreateUser(user entity.User) error {
	return r.db.Create(&user).Error
}
