package apihttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	api "oapi-generated/generated-api"
	"oapi-generated/internal/domain"
	"oapi-generated/internal/middleware"
)

func TestHandlersRequireAuthorization(t *testing.T) {
	app := newTestApp(&fakeUserService{})

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}

	var body api.Error
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if body.Message != "unauthorized" {
		t.Fatalf("message = %q, want unauthorized", body.Message)
	}
}

func TestHandlersListUsers(t *testing.T) {
	users := []domain.User{
		newHandlerTestUser("a@example.com", "Alice"),
		newHandlerTestUser("b@example.com", "Bob"),
	}
	svc := &fakeUserService{listUsers: users, total: 2}
	app := newTestApp(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body listUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if body.Total != 2 || len(body.Items) != 2 {
		t.Fatalf("total = %d len(items) = %d, want 2 and 2", body.Total, len(body.Items))
	}
	if svc.gotLimit != 20 || svc.gotOffset != 0 {
		t.Fatalf("limit = %d offset = %d, want 20 and 0", svc.gotLimit, svc.gotOffset)
	}
}

func TestHandlersListUsersForwardsPagination(t *testing.T) {
	svc := &fakeUserService{}
	app := newTestApp(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/users?limit=1&offset=1", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if svc.gotLimit != 1 || svc.gotOffset != 1 {
		t.Fatalf("limit = %d offset = %d, want 1 and 1", svc.gotLimit, svc.gotOffset)
	}
}

func TestHandlersGetUserNotFound(t *testing.T) {
	svc := &fakeUserService{getErr: domain.ErrNotFound}
	app := newTestApp(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+uuid.NewString(), nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	var body api.Error
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if body.Message != "not found" {
		t.Fatalf("message = %q, want not found", body.Message)
	}
}

func TestHandlersGetUserBadUUID(t *testing.T) {
	app := newTestApp(&fakeUserService{})

	req := httptest.NewRequest(http.MethodGet, "/v1/users/not-a-uuid", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

type fakeUserService struct {
	listUsers []domain.User
	total     int
	listErr   error
	getUser   domain.User
	getErr    error

	gotLimit  int
	gotOffset int
}

func (s *fakeUserService) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	s.gotLimit = limit
	s.gotOffset = offset
	if s.listErr != nil {
		return nil, 0, s.listErr
	}
	return s.listUsers, s.total, nil
}

func (s *fakeUserService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if s.getErr != nil {
		return domain.User{}, s.getErr
	}
	if s.getUser.ID == uuid.Nil {
		return domain.User{}, errors.New("unexpected get")
	}
	return s.getUser, nil
}

func newTestApp(svc *fakeUserService) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	api.RegisterHandlersWithOptions(app, NewHandlers(svc), api.FiberServerOptions{
		BaseURL: "/v1",
		HandlerMiddlewares: []api.HandlerMiddlewareFunc{
			middleware.BearerAuth(),
		},
	})
	return app
}

func newHandlerTestUser(email, name string) domain.User {
	return domain.User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC),
	}
}
