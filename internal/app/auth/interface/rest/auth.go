package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/reodora-backend/internal/app/auth/usecase"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecaseItf
}

func NewAuthHandler(router fiber.Router, authUsecase usecase.AuthUsecaseItf) {
	handler := AuthHandler{authUsecase}

	router.Post("/auth/login", handler.Login)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
