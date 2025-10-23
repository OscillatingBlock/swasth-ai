package auth

import (
	"context"
	"swasthAI/internal/auth/models"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *models.OTP) error
	FindByPhone(phone string) (*models.OTP, error)
	IncrementAttempts(phone string) error
	Delete(ctx context.Context, phone string) error
}
