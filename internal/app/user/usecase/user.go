package usecase

import "github.com/jevvonn/reodora-backend/internal/app/user/repository"

type UserUsecaseItf interface{}

type UserUsecase struct {
	userRepo repository.UserPostgreSQLItf
}

func NewUserUsecase(userRepo repository.UserPostgreSQLItf) UserUsecaseItf {
	return &UserUsecase{userRepo}
}
