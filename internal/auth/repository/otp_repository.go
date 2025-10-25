package repository

import (
	"context"
	"time"
	"github.com/uptrace/bun"
	"swasthAI/internal/auth/models"
)

type OTPRepository struct {
	db *bun.DB
}

func NewOTPRepository(db *bun.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (o *OTPRepository) Create(ctx context.Context, otp *models.OTP) error {
	return nil
}

func (r *OTPRepository) FindByPhone(ctx context.Context, phone string) (models.OTP, error) {
	return models.OTP{}, nil
}

func (r *OTPRepository) IncrementAttempts(ctx context.Context, phone string) error {
	return nil
}

func (r *OTPRepository) Delete(ctx context.Context, phone string) error {
	return nil
}

func (r *OTPRepository) CountRecent(ctx context.Context, phone string, duration time.Time) (int, error) {
	return 0, nil
}
