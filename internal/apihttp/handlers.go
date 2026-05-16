package apihttp

import (
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"

	api "oapi-generated/generated-api"
	"oapi-generated/internal/service"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

type Handlers struct {
	users service.UserService
}

type listUsersResponse struct {
	Total int        `json:"total"`
	Items []api.User `json:"items"`
}

func NewHandlers(users service.UserService) *Handlers {
	return &Handlers{users: users}
}

func (h *Handlers) ListUsers(c *fiber.Ctx, params api.ListUsersParams) error {
	limit := defaultListLimit
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	users, total, err := h.users.List(c.UserContext(), limit, offset)
	if err != nil {
		return writeError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(listUsersResponse{
		Total: total,
		Items: toAPIUsers(users),
	})
}

func (h *Handlers) GetUser(c *fiber.Ctx, id openapi_types.UUID) error {
	user, err := h.users.Get(c.UserContext(), id)
	if err != nil {
		return writeError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(toAPIUser(user))
}

var _ api.ServerInterface = (*Handlers)(nil)
