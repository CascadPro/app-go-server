package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"cascade/config"
	"cascade/pkg/logger"

	"github.com/google/uuid"
)

const (
	RedisSessionFolder   string = "cascade__session:"
	RedisCacheFolder     string = "cascade__cache:"
	RedisRateLimitFolder string = "cascade__rate_limit:"
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func GetSessionByID(sessionID string) (Session, error) {
	var session Session

	sessionString, query_err := config.R.DB.Get(config.R.Ctx, RedisSessionFolder+sessionID).Result()
	if query_err != nil {
		logger.Error("❌ Error during Redis session folder query!", query_err)
		return session, query_err
	}

	decode_err := json.Unmarshal([]byte(sessionString), &session)
	if decode_err != nil {
		logger.Error("❌ Error during Redis session string decoding!", decode_err)
		return session, decode_err
	}

	return session, nil
}

func CreateSession(r *http.Request, sessionID string, userID uuid.UUID, ttl time.Duration) error {
	session := Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}

	jsonString, err := json.Marshal(session)
	if err != nil {
		logger.Error("Error during JSON encoding!", err)
		return err
	}

	return config.R.DB.Set(config.R.Ctx, RedisSessionFolder+sessionID, jsonString, ttl).Err()
}

func DeleteSession(sessionIDs ...string) error {
	var keys []string

	for _, id := range sessionIDs {
		keys = append(keys, RedisSessionFolder+id)
	}

	return config.R.DB.Del(config.R.Ctx, keys...).Err()
}
