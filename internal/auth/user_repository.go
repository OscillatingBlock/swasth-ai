package auth

import (
	"context"
	"swasthAI/internal/auth/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}
