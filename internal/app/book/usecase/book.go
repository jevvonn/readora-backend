package usecase

import (
	"github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
)

type BookUsecaseItf interface{}

type BookUsecase struct {
	bookRepo repository.BookPostgreSQLItf
	log      logger.LoggerItf
}

func NewBookUsecase(userRepo repository.BookPostgreSQLItf, log logger.LoggerItf) BookUsecaseItf {
	return &BookUsecase{userRepo, log}
}
