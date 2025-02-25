package usecase

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/helper"
	authRepo "github.com/jevvonn/readora-backend/internal/app/auth/repository"
	userRepo "github.com/jevvonn/readora-backend/internal/app/user/repository"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/jwt"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/mailer"
	"github.com/jevvonn/readora-backend/internal/models"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(ctx *fiber.Ctx, req dto.RegisterRequest) error
	Login(ctx *fiber.Ctx, req dto.LoginRequest) (dto.LoginResponse, error)
	Session(ctx *fiber.Ctx) (dto.SessionResponse, error)
	SendRegisterOTP(ctx *fiber.Ctx, email string) error
	CheckRegisterOTP(ctx *fiber.Ctx, email, otp string) error
}

type AuthUsecase struct {
	userRepo userRepo.UserPostgreSQLItf
	authRepo authRepo.AuthRepositoryItf
	mailer   mailer.MailerItf
	log      logger.LoggerItf
}

func NewAuthUsecase(
	userRepo userRepo.UserPostgreSQLItf,
	authRepo authRepo.AuthRepositoryItf,
	log logger.LoggerItf,
	mailer mailer.MailerItf,
) AuthUsecaseItf {
	return &AuthUsecase{userRepo, authRepo, mailer, log}
}

func (u *AuthUsecase) Register(ctx *fiber.Ctx, req dto.RegisterRequest) error {
	log := "[AuthUsecase][Register]"

	// Check if username already exists
	user, err := u.userRepo.GetUserByEmailOrUsername(req.Email, req.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.Error(log, err)
		return err
	}

	if user.ID != uuid.Nil {
		return errors.New("user with this username or email already exists")
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

	// Send OTP to email
	err = u.SendRegisterOTP(ctx, req.Email)
	if err != nil {
		u.log.Error(log, err)
	}

	u.log.Info(log, "User created successfully")
	return nil
}

func (u *AuthUsecase) Login(ctx *fiber.Ctx, req dto.LoginRequest) (dto.LoginResponse, error) {
	log := "[AuthUsecase][Login]"

	// Check if username exists
	user, err := u.userRepo.GetUserByEmailOrUsername(req.Username, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.Error(log, err)
		return dto.LoginResponse{}, err
	}

	if user.ID == uuid.Nil {
		return dto.LoginResponse{}, models.ErrEmailOrUsernameExists
	}

	if !user.EmailVerified {
		return dto.LoginResponse{}, models.ErrEmailNotVerified
	}

	// Check password
	if !helper.VerifyPassword(req.Password, user.Password) {
		return dto.LoginResponse{}, models.ErrEmailOrUsernameExists
	}

	// Create Jwt token
	token, err := jwt.CreateAuthToken(user.ID.String(), user.Username)

	if err != nil {
		u.log.Error(log, err)
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		UserId: user.ID.String(),
		Token:  token,
	}, nil
}

func (u *AuthUsecase) Session(ctx *fiber.Ctx) (dto.SessionResponse, error) {
	log := "[AuthUsecase][Session]"
	userId := ctx.Locals("userId").(string)

	uuidUser, err := uuid.Parse(userId)
	if err != nil {
		u.log.Error(log, err)
		return dto.SessionResponse{}, err
	}

	user, err := u.userRepo.GetSpecificUser(entity.User{
		ID: uuidUser,
	})
	if err != nil {
		u.log.Error(log, err)
		return dto.SessionResponse{}, err
	}

	return dto.SessionResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (u *AuthUsecase) SendRegisterOTP(ctx *fiber.Ctx, email string) error {
	log := "[AuthUsecase][SendRegisterOTP]"

	// Check user exists
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Email: email,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.Error(log, err)
		return err
	}

	if user.ID == uuid.Nil {
		return models.ErrEmailNotExists
	}

	// Check if email already verified
	if user.EmailVerified {
		return models.ErrEmailAlreadyVerified
	}

	// Generate OTP
	otp := helper.RandomNumber(6)

	// Save OTP to Redis
	err = u.authRepo.SetRegisterOTP(ctx.Context(), email, otp)
	if err != nil {
		u.log.Error(log, err)
		return err
	}

	// Send OTP to email
	go func() {
		err = u.mailer.Send([]string{email}, "Register OTP", "Your OTP is "+otp)
		if err != nil {
			u.log.Error(log, err)
		}
	}()

	return nil
}

func (u *AuthUsecase) CheckRegisterOTP(ctx *fiber.Ctx, email, otp string) error {
	log := "[AuthUsecase][CheckRegisterOTP]"

	// Check user exists
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Email: email,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.Error(log, err)
		return err
	}

	if user.ID == uuid.Nil {
		return models.ErrEmailNotExists
	}

	// Check if email already verified
	if user.EmailVerified {
		return models.ErrEmailAlreadyVerified
	}

	// Get OTP from Redis
	savedOTP, err := u.authRepo.GetRegisterOTP(ctx.Context(), email)
	if err != nil {
		u.log.Error(log, err)
		return err
	}

	// Check OTP
	if savedOTP != otp {
		return models.ErrInvalidOTP
	}

	// Update user email verified
	err = u.userRepo.UpdateUser(entity.User{
		ID:            user.ID,
		EmailVerified: true,
	})
	if err != nil {
		u.log.Error(log, err)
		return err
	}

	// Delete OTP from Redis
	err = u.authRepo.DeleteRegisterOTP(ctx.Context(), email)
	if err != nil {
		u.log.Error(log, err)
	}

	return nil
}
