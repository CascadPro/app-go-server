package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

func GetSessionByID(userID uuid.UUID, sessionID string) (Session, error) {
	var session Session

	key := fmt.Sprintf("%s%s:%s", RedisSessionFolder, userID, sessionID)

	sessionString, queryErr := config.R.DB.Get(config.R.Ctx, key).Result()
	if queryErr != nil {
		logger.Error("❌ Error during Redis session folder query!", queryErr)
		return session, queryErr
	}

	decodeErr := json.Unmarshal([]byte(sessionString), &session)
	if decodeErr != nil {
		logger.Error("❌ Error during Redis session string decoding!", decodeErr)
		return session, decodeErr
	}

	return session, nil
}

func CreateSession(r *http.Request, sessionID string, userID uuid.UUID, ttl time.Duration) error {
	session := Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}

	metadata, metadataErr := getSessionMetadata(r)
	if metadataErr != nil {
		logger.Error("Failed to get session metadata!", metadataErr)
		return metadataErr
	}

	session.Metadata = *metadata

	jsonString, err := json.Marshal(session)
	if err != nil {
		logger.Error("Error during JSON encoding!", err)
		return err
	}

	key := fmt.Sprintf("%s%s:%s", RedisSessionFolder, userID, sessionID)

	return config.R.DB.Set(config.R.Ctx, key, jsonString, ttl).Err()
}

func DeleteUserSession(userID string, sessionIDs ...string) error {
	var keys []string

	key := fmt.Sprintf("%s%s:", RedisSessionFolder, userID)

	for _, id := range sessionIDs {
		keys = append(keys, key+id)
	}

	return config.R.DB.Del(config.R.Ctx, keys...).Err()
}
