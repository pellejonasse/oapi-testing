package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	api "oapi-generated/generated-api"
	"oapi-generated/internal/apihttp"
	"oapi-generated/internal/domain"
	"oapi-generated/internal/middleware"
	"oapi-generated/internal/service"
)

func main() {
	app := fiber.New(fiber.Config{ErrorHandler: apihttp.ErrorHandler})

	users := service.NewInMemoryUserService(seedUsers()...)
	handlers := apihttp.NewHandlers(users)

	api.RegisterHandlersWithOptions(app, handlers, api.FiberServerOptions{
		BaseURL: "/v1",
		HandlerMiddlewares: []api.HandlerMiddlewareFunc{
			middleware.BearerAuth(),
		},
	})

	log.Fatal(app.Listen(":8080"))
}

func seedUsers() []domain.User {
	now := time.Now().UTC()
	return []domain.User{
		{
			ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Email:     "alice@example.com",
			Name:      "Alice Example",
			CreatedAt: now,
		},
		{
			ID:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Email:     "bob@example.com",
			Name:      "Bob Example",
			CreatedAt: now,
		},
	}
}
