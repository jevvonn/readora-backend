package postgresql

import (
	"github.com/google/uuid"
	"github.com/jevvonn/reodora-backend/helper"
	"github.com/jevvonn/reodora-backend/internal/domain/entity"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	hashedPassword, err := helper.HashPassword("password")
	if err != nil {
		panic(err)
	}

	adminAccount := entity.User{
		ID:            uuid.New(),
		Name:          "Admin",
		Username:      "admin.ganteng",
		Email:         "admin@gmail.com",
		Password:      hashedPassword,
		Role:          "ADMIN",
		EmailVerified: true,
	}

	userAccount := entity.User{
		ID:            uuid.New(),
		Name:          "User",
		Username:      "user.ganteng",
		Email:         "user@gmail.com",
		Password:      hashedPassword,
		Role:          "USER",
		EmailVerified: true,
	}

	err = db.Create(&adminAccount).Error
	if err != nil {
		panic(err)
	}

	err = db.Create(&userAccount).Error
	if err != nil {
		panic(err)
	}
}
