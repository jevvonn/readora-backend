package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/auth/usecase"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/middleware"
	"github.com/jevvonn/readora-backend/internal/models"
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
	router.Post("/auth/otp", handler.SendRegisterOTP)
	router.Post("/auth/otp/check", handler.CheckRegisterOTP)

	router.Get("/auth/session", middleware.Authenticated, handler.Session)
}

// @Summary      Login User
// @Description  Login User
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.LoginRequest true "Login Request"
// @Success      200  object   models.JSONResponseModel{data=dto.LoginRequest,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Router       /api/auth/login [post]
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

// @Summary      Register User
// @Description  Register User
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.RegisterRequest true "Register Request"
// @Success      200  object   models.JSONResponseModel{data=dto.RegisterRequest,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Router       /api/auth/register [post]
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

// @Summary      Get Session User Data
// @Description  Get Session User Data
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  object   models.JSONResponseModel{data=dto.SessionResponse,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/auth/session [get]
func (h *AuthHandler) Session(ctx *fiber.Ctx) error {
	log := "[AuthHandler][Session]"

	res, err := h.authUsecase.Session(ctx)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	return h.response.SetData(res).Success(ctx, "Session Retrieved Successfully")
}

func (h *AuthHandler) SendRegisterOTP(ctx *fiber.Ctx) error {
	log := "[AuthHandler][SendRegisterOTP]"

	var req dto.SendRegisterOTPRequest
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

	err = h.authUsecase.SendRegisterOTP(ctx, req.Email)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	return h.response.Success(ctx, "OTP Sent Successfully")
}

func (h *AuthHandler) CheckRegisterOTP(ctx *fiber.Ctx) error {
	log := "[AuthHandler][CheckRegisterOTP]"

	var req dto.CheckRegisterOTPRequest
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

	err = h.authUsecase.CheckRegisterOTP(ctx, req.Email, req.OTP)
	if err != nil {
		h.log.Error(log, err)
		return h.response.BadRequest(ctx, err, nil)
	}

	return h.response.Success(ctx, "Email Verified Successfully")
}
