package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/infra/jwt"
)

func Authenticated(ctx *fiber.Ctx) error {
	headers := ctx.Get("Authorization")

	if headers == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})

	}

	tokenString := strings.Replace(headers, "Bearer ", "", 1)
	claims, err := jwt.ParseAuthToken(tokenString)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	exp := claims["exp"].(float64)
	expiredDate := time.Unix(int64(exp), 0)

	if expiredDate.Before(time.Now()) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token expired",
		})
	}

	ctx.Locals("userId", claims["sub"])
	ctx.Locals("username", claims["username"])
	ctx.Locals("role", claims["role"])

	return ctx.Next()
}
