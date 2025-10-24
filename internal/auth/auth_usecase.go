package auth

import (
	"context"
	"swasthAI/internal/auth/models"
)

type AuthUsecase interface {
	SendOTP(ctx context.Context, phone string) error
	VerifyOTP(ctx context.Context, phone, otp string) (*models.UserWithToken, bool, error)
	RegisterUser(ctx context.Context, input *models.RegisterUserInput) (*models.UserWithToken, error)
	GetUserByID(ctx context.Context, token string) (*models.User, error)
	UpdateProfile(ctx context.Context, input *models.UpdateProfileInput) (*models.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.Tokens, error)
	ResendOTP(ctx context.Context, phone string) error
}
