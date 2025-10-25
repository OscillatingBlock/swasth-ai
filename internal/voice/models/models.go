package voice

import (
	"time"
)

type WSVoiceMessage struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Language string `json:"language"`
	Mode     string `json:"mode"`
}

type TranscribeResponse struct {
	Transcription string    `json:"transcription"`
	Language      string    `json:"language"`
	Confidence    float64   `json:"confidence"`
	Duration      time.Time `json:"duration"` // seconds
}

type AnalyzeRequest struct {
	Query    string  `json:"query" validate:"required"`
	Language string  `json:"language" validate:"required,oneof=hi en ta te bn mr"`
	Mode     *string `json:"mode,omitempty"` // optional: "offline" or "online"
}

type AnalyzeResponse struct {
	Advice            string    `json:"advice"`
	Severity          string    `json:"severity"` // "low", "moderate", "high"
	FollowupQuestions []string  `json:"followup_questions"`
	DoctorReferral    bool      `json:"doctor_referral"`
	Language          string    `json:"language"`
	GeneratedAt       time.Time `json:"generated_at"`
}
