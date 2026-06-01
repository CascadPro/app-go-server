package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeRegister    TokenType = "REGISTER"
	TokenTypeEmailVerify TokenType = "EMAIL_VERIFY"
)

type Token struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;not null" json:"id"`

	Token string    `gorm:"type:varchar(255);uniqueIndex:idx_user_token;type:varchar(255);not null" json:"token"`
	Type  TokenType `gorm:"type:varchar(255);not null" json:"type"`

	User   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_token;not null" json:"user_id,omitempty"`

	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at,omitempty"`
}
