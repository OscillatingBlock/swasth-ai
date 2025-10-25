package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"swasthAI/config"
	"swasthAI/internal/auth/models"
	appErrors "swasthAI/pkg/errors"
	"time"
)

// func GenerateJWTToken(user *models.User, jwtConfig config.JWT) (string, string, error) {
// 	return "access-token", "refresh-token", nil
// }

// GenerateJWTToken generates access and refresh tokens for a user
func GenerateJWTToken(user *models.User, jwtConfig config.JWT) (string, string, error) {
	now := time.Now()

	// Access token
	accessClaims := JWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(jwtConfig.ExpiresIn) * time.Second)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (longer expiry, e.g., 7 days)
	refreshClaims := JWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)), // 7 days
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateRefreshToken(refreshToken string, jwtConfig config.JWT) (*models.User, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})
	if err != nil {
		return nil, appErrors.ErrInvalidJWTToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, appErrors.ErrInvalidJWTToken
	}

	// Optionally check expiry explicitly
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, appErrors.ErrJWTExpired
	}

	user := &models.User{ID: claims.ID}
	return user, nil
}

func ValidateToken(token string, jwtConfig config.JWT) (*models.User, error) {
	return nil, nil
}

type JWTClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}
