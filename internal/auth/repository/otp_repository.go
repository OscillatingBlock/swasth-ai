package repository

import (
	"context"
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
