package sessions

import "github.com/google/uuid"

type Session struct {
	UserID    uuid.UUID       `json:"user_id"`
	ExpiresAt int64           `json:"expires_at"`
	Metadata  SessionMetadata `json:"metadata"`
}

type SessionMetadata struct {
	IP       string   `json:"ip"`
	Location Location `json:"location"`
	Device   Device   `json:"device"`
}

type Location struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"latitude"`
	Lng     float64 `json:"longitude"`
}

type Device struct {
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Model   string `json:"model"`
	Type    string `json:"type"`
}
