package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type AuthRepositoryItf interface {
	SetRegisterOTP(ctx context.Context, email string, otp string) error
	GetRegisterOTP(ctx context.Context, email string) (string, error)
	DeleteRegisterOTP(ctx context.Context, email string) error
}

type AuthRepository struct {
	rdb *redis.Client
	log logger.LoggerItf
}

func NewAuthRepository(rdb *redis.Client, log logger.LoggerItf) AuthRepositoryItf {
	return &AuthRepository{rdb, log}
}

func (r *AuthRepository) SetRegisterOTP(ctx context.Context, email string, otp string) error {
	key := fmt.Sprintf("auth:register:otp:%s", email)
	err := r.rdb.SetEx(ctx, key, otp, time.Minute*10).Err()

	if err != nil {
		r.log.Error("[AuthRepository][SendRegisterOTP]", err)
	}

	return err
}

func (r *AuthRepository) GetRegisterOTP(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("auth:register:otp:%s", email)
	cmd := r.rdb.Get(ctx, key)

	err := cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.log.Error("[AuthRepository][SendRegisterOTP]", err)
			return "", models.ErrInvalidOTP
		}

		r.log.Error("[AuthRepository][SendRegisterOTP]", err)
		return "", err
	}

	otp, err := cmd.Result()
	if err != nil {
		r.log.Error("[AuthRepository][SendRegisterOTP]", err)
		return "", err
	}

	return otp, nil
}

func (r *AuthRepository) DeleteRegisterOTP(ctx context.Context, email string) error {
	key := fmt.Sprintf("auth:register:otp:%s", email)
	cmd := r.rdb.Del(ctx, key)

	err := cmd.Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		r.log.Error("[AuthRepository][SendRegisterOTP]", err)
		return err
	}

	return nil
}
