package voice

import (
	"context"

	"swasthAI/internal/voice/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type VoiceUseCase interface {
	StartSession(ctx context.Context, req *models.StartSessionRequest, UserID uuid.UUID) error
	HandleClientWebSocket(conn *websocket.Conn, sessionID uuid.UUID)
	EndSession(ctx context.Context, sessionID uuid.UUID) error
}
