package repository

import (
	"context"
	"sync"

	"swasthAI/internal/voice/models"
	"swasthAI/pkg/domain_errors"
)

type InMemorySessionRepository struct {
	sessions sync.Map
	mu       sync.RWMutex
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: sync.Map{},
		mu:       sync.RWMutex{},
	}
}

func (sr *InMemorySessionRepository) CreateSession(ctx context.Context, session *models.VoiceSession) error {
	sr.mu.Lock()
	sr.sessions.Store(session.SessionID, session)
	sr.mu.Unlock()
	return nil
}

func (sr *InMemorySessionRepository) GetSession(ctx context.Context, id string) (*models.VoiceSession, error) {
	sr.mu.RLock()
	session, ok := sr.sessions.Load(id)
	if !ok {
		return nil, domain_errors.ErrSessionNotFound
	}
	sr.mu.RUnlock()
	return session.(*models.VoiceSession), nil
}

func (sr *InMemorySessionRepository) UpdateSession(ctx context.Context, session *models.VoiceSession) error {
	sr.mu.Lock()
	sr.sessions.Store(session.SessionID, session)
	sr.mu.Unlock()
	return nil
}

func (sr *InMemorySessionRepository) DeleteSession(ctx context.Context, id string) error {
	sr.mu.Lock()
	sr.sessions.Delete(id)
	sr.mu.Unlock()
	return nil
}

func (sr *InMemorySessionRepository) ListActiveSessions(ctx context.Context) ([]models.VoiceSession, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	var sessions []models.VoiceSession
	for _, v := range sr.sessions.Range {
		sessions = append(sessions, *v.(*models.VoiceSession))
	}
	return sessions, nil
}
