package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeRegister    TokenType = "register"
	TokenTypeEmailVerify TokenType = "email_verify"
)

func (ct *TokenType) Scan(value any) error {
	*ct = TokenType(value.([]byte))
	return nil
}

func (ct TokenType) Value() (driver.Value, error) {
	return string(ct), nil
}

type Token struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`

	Token string    `gorm:"uniqueIndex:idx_user_token" json:"token"`
	Type  TokenType `gorm:"type:token_type" json:"type"`

	User   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_token" json:"user_id"`

	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
}
