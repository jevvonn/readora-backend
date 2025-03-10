package postgresql

import (
	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/helper"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
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
		Username:      "adminganteng",
		Email:         "admin@gmail.com",
		Password:      hashedPassword,
		Role:          constant.RoleAdmin,
		EmailVerified: true,
	}

	userAccount := entity.User{
		ID:            uuid.New(),
		Name:          "User",
		Username:      "userganteng",
		Email:         "user@gmail.com",
		Password:      hashedPassword,
		Role:          constant.RoleUser,
		EmailVerified: true,
	}

	genres := []entity.Genre{
		{Name: "Fiction"},
		{Name: "Non-Fiction"},
		{Name: "Fantasy"},
		{Name: "Science Fiction"},
		{Name: "Romance"},
		{Name: "Mystery"},
		{Name: "Horror"},
		{Name: "Thriller"},
		{Name: "Historical Fiction"},
		{Name: "Young Adult"},
	}

	err = db.Create(&genres).Error
	if err != nil {
		panic(err)
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
