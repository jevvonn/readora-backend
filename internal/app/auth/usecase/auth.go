package usecase

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jevvonn/reodora-backend/helper"
	"github.com/jevvonn/reodora-backend/internal/app/user/repository"
	"github.com/jevvonn/reodora-backend/internal/constant"
	"github.com/jevvonn/reodora-backend/internal/domain/dto"
	"github.com/jevvonn/reodora-backend/internal/domain/entity"
	"github.com/jevvonn/reodora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(ctx *fiber.Ctx, req dto.RegisterRequest) error
}

type AuthUsecase struct {
	userRepo repository.UserPostgreSQLItf
	log      logger.LoggerItf
}

func NewAuthUsecase(userRepo repository.UserPostgreSQLItf, log logger.LoggerItf) AuthUsecaseItf {
	return &AuthUsecase{userRepo, log}
}

func (u *AuthUsecase) Register(ctx *fiber.Ctx, req dto.RegisterRequest) error {
	log := "[AuthUsecase][Register]"

	// Check if username already exists
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Username: req.Username,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.Error(log, err)
		return err
	}

	if user.ID != uuid.Nil {
		return errors.New("user with this username already exists")
	}

	// Check if email already exists
	user, err = u.userRepo.GetSpecificUser(entity.User{
		Email: req.Email,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.Error(log, err)
		return err
	}

	if user.ID != uuid.Nil {
		return errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		u.log.Error(log, err)
		return err
	}

	// Create user
	user = entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     constant.RoleUser,
	}
	err = u.userRepo.CreateUser(user)

	if err != nil {
		u.log.Error(log, err)
		return err
	}

	u.log.Info(log, "User created successfully")
	return nil
}
