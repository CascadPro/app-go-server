package authutils

import (
	"fmt"
	"time"

	"cascade/internal/models"
	"cascade/pkg/logger"
	"cascade/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	UserID    uuid.UUID       `json:"user_id"`
	Role      models.UserRole `json:"role"`
	SessionID string          `json:"session_id"`
	ExpiresAt int64           `json:"exp"`

	jwt.RegisteredClaims
}

const (
	RefreshTokenLifetime time.Duration = time.Hour * 24 * 30
	AccessTokenLifetime  time.Duration = time.Minute * 15
)

func GenerateToken(method jwt.SigningMethod, claims jwt.Claims, cfg *utils.Config) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, cfg *utils.Config) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неожиданный метод подписи: %v", t.Header["alg"])
		}

		return []byte(cfg.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("Ключ невалиден")
	}

	claims, ok := token.Claims.(*TokenPayload)
	if !ok {
		return nil, fmt.Errorf("Не удалось получить payload")
	}

	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, fmt.Errorf("Срок действия ключа закончился")
	}

	return claims, nil
}

func IssueTokens(userID uuid.UUID, role models.UserRole, sessionID string, cfg *utils.Config) (string, string) {
	accessToken, accessTokenErr := GenerateToken(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"role":       role,
		"session_id": sessionID,
		"exp":        time.Now().Add(AccessTokenLifetime).Unix(),
	}, cfg)
	if accessTokenErr != nil {
		accessToken = ""
		logger.Error("Error during access token generation!", accessTokenErr)
	}

	refreshToken, refreshTokenErr := GenerateToken(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"role":       role,
		"session_id": sessionID,
		"exp":        time.Now().Add(RefreshTokenLifetime).Unix(),
	}, cfg)
	if refreshTokenErr != nil {
		refreshToken = ""
		logger.Error("Error during refresh token generation!", refreshTokenErr)
	}

	return accessToken, refreshToken
}
