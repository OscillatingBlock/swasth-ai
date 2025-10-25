package middleware

import (
	"context"
	"strings"
	"swasthAI/config"
	"swasthAI/internal/auth/usecase"
	appErrors "swasthAI/pkg/errors"
	"swasthAI/pkg/http_errors"
	"swasthAI/pkg/utils"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (mw *MiddlewareManager) AuthJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerHeader := c.Request().Header.Get("Authorization")
		if bearerHeader != "" {
			headerParts := strings.Split(bearerHeader, " ")
			if len(headerParts) != 2 {
				mw.Logger.Error("len headerParts != 2")
				return http_errors.Send(c, appErrors.ErrInvalidInput)
			}

			tokenString := headerParts[1]
			err := mw.ValidateJWTToken(tokenString, mw.AuthUC, mw.Cfg, c)
			if err != nil {
				return err
			}
			return next(c)
		}
		return next(c)
	}
}

func (mw MiddlewareManager) ValidateJWTToken(tokenString string, authUC usecase.AuthUsecase, cfg config.Config, c echo.Context) error {
	fmt.Printf("ValidateJWTToken %s\n", tokenString)
	if tokenString == "" {
		return appErrors.ErrInvalidJWTToken
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.JWT.Secret)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return appErrors.ErrInvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("claims", claims)
		userID, ok := claims["id"].(string)
		if !ok {
			return appErrors.ErrJWTInvalidClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		jwtClaims := &utils.JWTClaims{
			ID: userUUID,
		}
		fmt.Printf("set claims %v\n", jwtClaims)
		ctx := context.WithValue(c.Request().Context(), "claims", jwtClaims)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}

func (mw *MiddlewareManager) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)

		req := c.Request()
		res := c.Response()

		mw.Logger.Infof(
			"[HTTP] method=%s path=%s status=%d latency=%s ip=%s user_agent=%q request_id=%s error=%v",
			req.Method,
			req.URL.Path,
			res.Status,
			time.Since(start),
			c.RealIP(),
			req.UserAgent(),
			req.Header.Get("X-Request-ID"),
			err,
		)
		return err
	}
}
