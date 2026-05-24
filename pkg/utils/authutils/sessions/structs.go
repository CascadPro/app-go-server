package sessions

import "github.com/google/uuid"

type Session struct {
	UserID    uuid.UUID       `json:"user_id"`
	ExpiresAt int64           `json:"expires_at"`
}

