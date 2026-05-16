package apihttp

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	api "oapi-generated/generated-api"
	"oapi-generated/internal/domain"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	return writeError(c, err)
}

func writeError(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "internal error"

	var fiberErr *fiber.Error
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = fiber.StatusNotFound
		message = domain.ErrNotFound.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status = fiber.StatusUnauthorized
		message = domain.ErrUnauthorized.Error()
	case errors.As(err, &fiberErr):
		status = fiberErr.Code
		message = fiberErr.Message
	}

	return c.Status(status).JSON(api.Error{Message: message})
}
