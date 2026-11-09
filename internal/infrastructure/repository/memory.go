package repository

import (
	"context"
	"mt-bot/internal/domain"
	"mt-bot/internal/usecase"
	"sync"
)

type MemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[int64]domain.UserSession
}

func NewMemorySessionRepository() *MemorySessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[int64]domain.UserSession),
	}
}

func (m *MemorySessionRepository) Get(ctx context.Context, chatID int64) (domain.UserSession, error) {
	if err := ctx.Err(); err != nil {
		return domain.UserSession{}, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[chatID]
	if !ok {
		return domain.UserSession{}, usecase.ErrSessionNotFound
	}

	return session, nil

}

func (m *MemorySessionRepository) Save(
	ctx context.Context,
	chatID int64,
	session domain.UserSession,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[chatID] = session

	return nil
}

func (m *MemorySessionRepository) Delete(ctx context.Context, chatID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, chatID)

	return nil
}

var _ usecase.SessionRepository = (*MemorySessionRepository)(nil)
