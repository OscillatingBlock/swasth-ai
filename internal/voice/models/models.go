// internal/voice/models/models.go
package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type StartSessionRequest struct {
	Language    string `json:"language"`
	Model       string `json:"model"`
	SessionType string `json:"session_type"`
}

type StartSessionResponse struct {
	SessionID string `json:"session_id"`
	WSURL     string `json:"ws_url"`
}

type EndSessionRequest struct {
	SessionID string `json:"session_id"`
}

// Internal session state
type VoiceSession struct {
	SessionID string
	UserID    string
	Language  string
	Model     string
	CreatedAt time.Time
	ExpiresAt time.Time
	AiWSConn  *websocket.Conn
	Status    string
}

// WebSocket transport
type WSMessage struct {
	Type string `json:"type"`
}

// Client → Server
type TextMessageInput struct {
	Content string `json:"content,omitempty"`
}

type EndOfInput struct{}

// Server → Client
type Transcript struct {
	Text string `json:"text,omitempty"`
}

type AIText struct {
	Text string `json:"ai_text,omitempty"`
}

type EndOfResponse struct{}

type SessionStrore struct {
	sessions map[string]*VoiceSession
	mu       sync.RWMutex
}
