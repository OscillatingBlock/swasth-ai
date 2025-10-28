package voice

import (
	"context"
	"swasthAI/internal/voice/models"

	"github.com/google/uuid"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.VoiceSession) error
	GetSession(ctx context.Context, id uuid.UUID) (*models.VoiceSession, error)
	UpdateSession(ctx context.Context, session *models.VoiceSession) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	ListActiveSessions(ctx context.Context) ([]models.VoiceSession, error)
}
