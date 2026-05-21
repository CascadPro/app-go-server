package auth

import (
	"time"

	"cascade/internal/models"
	"cascade/pkg/logger"
	"cascade/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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
