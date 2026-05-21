package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"cascade/config"
	"cascade/pkg/logger"

	"github.com/google/uuid"
)

type Session struct {
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt int64     `json:"expires_at"`
}

const (
	RedisSessionFolder string = "cascade__session:"
	RedisCacheFolder   string = "cascade__cache:"
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func CreateSession(sessionID string, userID uuid.UUID, ttl time.Duration) error {
	session := Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}

	jsonString, err := json.Marshal(session)
	if err != nil {
		logger.Error("Error during JSON encoding!", err)
	}

	return config.Redis.Set(config.RedisCtx, RedisSessionFolder+sessionID, jsonString, ttl).Err()
}
