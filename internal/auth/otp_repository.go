package auth

import (
	"context"
	"swasthAI/internal/auth/models"
	"time"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *models.OTP) error
	FindByPhone(ctx context.Context, phone string) (models.OTP, error)
	IncrementAttempts(ctx context.Context, phone string) error
	Delete(ctx context.Context, phone string) error
	CountRecent(ctx context.Context, phone string, duration time.Time) (int, error)
}
