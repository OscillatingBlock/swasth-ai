package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"swasthAI/config"
	"swasthAI/internal/voice/models"
	"swasthAI/internal/voice/repository"
	"swasthAI/pkg/domain_errors"
	appErrors "swasthAI/pkg/errors"
	"swasthAI/pkg/logger"
	"swasthAI/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type VoiceUsecase struct {
	SessionRepo *repository.InMemorySessionRepository
	aiWSURL     string
	logger      *logger.Logger
	upgrader    websocket.Upgrader
	httpClient  *http.Client
	config      *config.Config
}

func NewVoiceUsecase(logger *logger.Logger, SessionRepo *repository.InMemorySessionRepository, aiWSURL string, httpClient *http.Client) *VoiceUsecase {
	return &VoiceUsecase{SessionRepo: SessionRepo, logger: logger, aiWSURL: aiWSURL, httpClient: httpClient}
}

func (u *VoiceUsecase) StartSession(ctx context.Context, req *models.StartSessionRequest, UserID uuid.UUID) (*models.StartSessionResponse, error) {
	claims, ok := ctx.Value("claims").(*utils.JWTClaims)
	if !ok {
		return nil, appErrors.ErrUnauthorized
	}
	userID := claims.ID

	sessionUUID := uuid.New()
	shortID := "vsn_" + sessionUUID.String()[:6]
	aiConn, _, err := websocket.DefaultDialer.Dial(u.aiWSURL+"?session_id="+shortID, nil)
	if err != nil {
		return nil, domain_errors.ErrAIConnectionFailed
	}

	session := &models.VoiceSession{
		SessionID: shortID,
		UserID:    userID.String(),
		Language:  req.Language,
		Model:     req.Model,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(10 * 60),
		AiWSConn:  aiConn,
		Status:    "active",
	}

	err = u.SessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return &models.StartSessionResponse{
		SessionID: shortID,
		WSURL:     u.aiWSURL + "/api/v1/voice/session/" + shortID + "/ws",
	}, nil
}

func (u *VoiceUsecase) HandleClientWebSocket(ctx context.Context, clientConn *websocket.Conn, sessionID string) {
	session, err := u.SessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return
	}
	toAI := make(chan []byte, 100)
	fromAI := make(chan []byte, 100)
	done := make(chan struct{})

	go u.relayFromAI(session.AiWSConn, clientConn, fromAI, done)

	for {
		msgType, data, err := clientConn.ReadMessage()
		if err != nil {
			break
		}

		switch msgType {
		case websocket.BinaryMessage:
			select {
			case toAI <- data:
			case <-ctx.Done():
				return
			}

		case websocket.TextMessage:
			var msg *models.WSMessage
			if json.Unmarshal(data, &msg) != nil {
				continue
			}

			switch msg.Type {
			case "end_of_input":
				session.AiWSConn.WriteJSON(map[string]any{"type": "end_of_input"})
			case "text_message":
				var input *models.TextMessageInput
				if json.Unmarshal(data, &input) == nil {
					session.AiWSConn.WriteJSON(map[string]any{
						"type": "text_message", "content": input.Content,
					})
				}
			}

		}

	}
	close(done)
	clientConn.Close()
	session.AiWSConn.Close()
	u.SessionRepo.DeleteSession(ctx, sessionID)
}

func (u *VoiceUsecase) relayFromAI(aiConn, clientConn *websocket.Conn, fromAI chan []byte, done chan struct{}) {
	defer close(fromAI)
	for {
		msgType, data, err := aiConn.ReadMessage()
		if err != nil {
			break
		}

		switch msgType {
		case websocket.BinaryMessage:
			// ai_audio
			clientConn.WriteMessage(websocket.BinaryMessage, data)
		case websocket.TextMessage:
			var msg models.WSMessage
			if json.Unmarshal(data, &msg) == nil {
				switch msg.Type {
				case "partial_transcript", "final_transcript", "ai_text":
					clientConn.WriteJSON(map[string]any{"type": msg.Type, "text": "..."})
				case "end_of_response":
					clientConn.WriteJSON(map[string]any{"type": "end_of_response"})
				}
			}
		}
	}
	select {
	case <-done:
	default:
		close(done)
	}
}
