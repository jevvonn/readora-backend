package repository

import(
    "gorm.io/gorm"
)

type UserPostgreSQLItf interface {}

type UserPostgreSQL struct {
    db *gorm.DB
}

func NewUserPostgreSQL(db *gorm.DB) UserPostgreSQLItf {
    return &UserPostgreSQL{db}
}
