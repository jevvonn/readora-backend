package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/reodora-backend/internal/app/auth/usecase"
	"github.com/jevvonn/reodora-backend/internal/domain/dto"
	"github.com/jevvonn/reodora-backend/internal/infra/logger"
	"github.com/jevvonn/reodora-backend/internal/infra/validator"
	"github.com/jevvonn/reodora-backend/internal/models"
)

type AuthHandler struct {
	router      fiber.Router
	authUsecase usecase.AuthUsecaseItf
	validator   validator.ValidationService
	log         logger.LoggerItf
	response    models.ResponseItf
}

func NewAuthHandler(
	router fiber.Router,
	authUsecase usecase.AuthUsecaseItf,
	validator validator.ValidationService,
	log logger.LoggerItf,
	response models.ResponseItf,
) {
	handler := AuthHandler{router, authUsecase, validator, log, response}

	router.Post("/auth/login", handler.Login)
	router.Post("/auth/register", handler.Register)
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	log := "[AuthHandler][Login]"

	var req dto.LoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	erorrsMap, err := h.validator.Validate(req)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, erorrsMap)
	}

	res, err := h.authUsecase.Login(ctx, req)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	return h.response.SetData(res).Success(ctx, "User Logged In Successfully")
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	log := "[AuthHandler][Register]"

	var req dto.RegisterRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	erorrsMap, err := h.validator.Validate(req)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, erorrsMap)
	}

	err = h.authUsecase.Register(ctx, req)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	return h.response.Success(ctx, "User Registered Successfully")
}
