package utils

import (
	"swasthAI/config"
	"swasthAI/internal/auth/models"

	"github.com/google/uuid"
)

func GenerateJWTToken(user *models.User, jwtConfig config.JWT) (string, string, error) {
	return "", "", nil
}

func ValidateRefreshToken(refreshToken string, jwtConfig config.JWT) (*models.User, error) {
	return nil, nil
}

func ValidateToken(token string, jwtConfig config.JWT) (*models.User, error) {
	return nil, nil
}

type JWTClaims struct {
	ID uuid.UUID
}
