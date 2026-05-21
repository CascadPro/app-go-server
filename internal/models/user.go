package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleRegular        UserRole = "REGULAR"
	RoleForeman        UserRole = "FOREMAN"
	RoleProjectManager UserRole = "PROJECT_MANAGER"
	RoleClerk          UserRole = "CLERK"
	RoleEngineer       UserRole = "ENGINEER"
	RoleDirector       UserRole = "DIRECTOR"
)

type User struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`

	Email        *string  `gorm:"unique;uniqueIndex" json:"email"`
	PasswordHash *string  `gorm:"column:password" json:"password_hash"`
	Role         UserRole `gorm:"not null" json:"role"`

	Token []Token `gorm:"foreignKey:UserID" json:"tokens"`

	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt int64     `gorm:"autoUpdateTime:unix" json:"updated_at"`
}
