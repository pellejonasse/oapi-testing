package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	api "oapi-generated/generated-api"
)

func BearerAuth() api.HandlerMiddlewareFunc {
	return func(c *fiber.Ctx, next fiber.Handler) error {
		header := c.Get(fiber.HeaderAuthorization)
		const prefix = "Bearer "

		if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) || strings.TrimSpace(header[len(prefix):]) == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		return next(c)
	}
}
