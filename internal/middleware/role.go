package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

func RequireRoles(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role := ctx.Locals("role").(string)

		if !slices.Contains(roles, role) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden resource",
			})

		}

		return ctx.Next()
	}
}
