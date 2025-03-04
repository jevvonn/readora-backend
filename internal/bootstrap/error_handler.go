package bootstrap

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/models"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	var e *fiber.Error
	if errors.As(err, &e) {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.JSONResponseModel{
			Message: e.Message,
		})
	}

	var valErr *validator.ValidationError
	if errors.As(err, &valErr) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(models.JSONResponseModel{
			Message: valErr.Message,
			Errors:  valErr.Errors,
		})
	}

	var errResp *errorpkg.ErrorResp
	if errors.As(err, &errResp) {
		return ctx.Status(errResp.StatusCode).JSON(models.JSONResponseModel{
			Message: errResp.Message,
		})
	}

	return nil
}
