package service

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"oapi-generated/internal/domain"
)

const (
	defaultUserLimit = 20
	maxUserLimit     = 100
)

type UserService interface {
	List(ctx context.Context, limit, offset int) ([]domain.User, int, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type InMemoryUserService struct {
	mu    sync.RWMutex
	users map[uuid.UUID]domain.User
	order []uuid.UUID
}

func NewInMemoryUserService(seed ...domain.User) *InMemoryUserService {
	s := &InMemoryUserService{
		users: make(map[uuid.UUID]domain.User, len(seed)),
		order: make([]uuid.UUID, 0, len(seed)),
	}

	for _, user := range seed {
		s.upsert(user)
	}

	return s
}

func (s *InMemoryUserService) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	limit = normalizeLimit(limit)
	if offset < 0 {
		offset = 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.order)
	if offset >= total {
		return []domain.User{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	items := make([]domain.User, 0, end-offset)
	for _, id := range s.order[offset:end] {
		items = append(items, s.users[id])
	}

	return items, total, nil
}

func (s *InMemoryUserService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	select {
	case <-ctx.Done():
		return domain.User{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}

	return user, nil
}

func (s *InMemoryUserService) upsert(user domain.User) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	if _, exists := s.users[user.ID]; !exists {
		s.order = append(s.order, user.ID)
	}
	s.users[user.ID] = user
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultUserLimit
	}
	if limit > maxUserLimit {
		return maxUserLimit
	}
	return limit
}
