package apihttp

import (
	api "oapi-generated/generated-api"
	"oapi-generated/internal/domain"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

func toAPIUser(user domain.User) api.User {
	return api.User{
		Id:        openapi_types.UUID(user.ID),
		Email:     openapi_types.Email(user.Email),
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func toAPIUsers(users []domain.User) []api.User {
	items := make([]api.User, 0, len(users))
	for _, user := range users {
		items = append(items, toAPIUser(user))
	}
	return items
}
