package usecase

import "github.com/jevvonn/reodora-backend/internal/app/user/repository"

type AuthUsecaseItf interface{}

type AuthUsecase struct {
	userRepo repository.UserPostgreSQLItf
}

func NewAuthUsecase(userRepo repository.UserPostgreSQLItf) AuthUsecaseItf {
	return &AuthUsecase{userRepo}
}
