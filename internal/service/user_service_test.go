package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"oapi-generated/internal/domain"
)

func TestInMemoryUserServiceList(t *testing.T) {
	users := []domain.User{
		newTestUser("a@example.com", "Alice"),
		newTestUser("b@example.com", "Bob"),
		newTestUser("c@example.com", "Carol"),
	}
	svc := NewInMemoryUserService(users...)

	items, total, err := svc.List(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if total != 3 {
		t.Fatalf("total = %d, want 3", total)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != users[1].ID {
		t.Fatalf("items[0].ID = %s, want %s", items[0].ID, users[1].ID)
	}
}

func TestInMemoryUserServiceListDefaultsAndBounds(t *testing.T) {
	users := []domain.User{newTestUser("a@example.com", "Alice")}
	svc := NewInMemoryUserService(users...)

	items, total, err := svc.List(context.Background(), 0, -10)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("total = %d len(items) = %d, want 1 and 1", total, len(items))
	}

	items, total, err = svc.List(context.Background(), 10, 99)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if total != 1 || len(items) != 0 {
		t.Fatalf("total = %d len(items) = %d, want 1 and 0", total, len(items))
	}
}

func TestInMemoryUserServiceGet(t *testing.T) {
	user := newTestUser("a@example.com", "Alice")
	svc := NewInMemoryUserService(user)

	got, err := svc.Get(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.ID != user.ID {
		t.Fatalf("got ID = %s, want %s", got.ID, user.ID)
	}
}

func TestInMemoryUserServiceGetNotFound(t *testing.T) {
	svc := NewInMemoryUserService()

	_, err := svc.Get(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func newTestUser(email, name string) domain.User {
	return domain.User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC),
	}
}
