package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/redis/go-redis/v9"
)

type AuthRepositoryItf interface {
	SetRegisterOTP(ctx context.Context, email string, otp string) error
	GetRegisterOTP(ctx context.Context, email string) (string, error)
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

	if cmd.Err() != nil {
		r.log.Error("[AuthRepository][SendRegisterOTP]", cmd.Err())
		return "", cmd.Err()
	}

	otp, err := cmd.Result()

	if err != nil {
		r.log.Error("[AuthRepository][SendRegisterOTP]", err)
		return "", err
	}

	return otp, nil
}
