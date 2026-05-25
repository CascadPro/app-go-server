package models

import (
	"cascade/pkg/logger"
	"cascade/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleForeman        UserRole = "FOREMAN"
	RoleProjectManager UserRole = "PROJECT_MANAGER"
	RoleClerk          UserRole = "CLERK"
	RoleEngineer       UserRole = "ENGINEER"
	RoleDirector       UserRole = "DIRECTOR"
	RoleRegular        UserRole = "REGULAR"
	RoleAdmin          UserRole = "ADMIN"
)

type User struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`

	Email        *string  `gorm:"type:varchar(255);unique;uniqueIndex" json:"email"`
	PasswordHash *string  `gorm:"column:password;type:varchar(255)" json:"password_hash,omitempty"`
	Role         UserRole `gorm:"type:varchar(255);not null" json:"role"`

	Token []Token `gorm:"foreignKey:UserID" json:"tokens,omitempty"`

	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime:unix" json:"updated_at,omitempty"`
}

func (u *User) HashPassword() error {
	password_hash, err := utils.GenerateHash(*u.PasswordHash)
	if err != nil {
		return err
	}

	u.PasswordHash = &password_hash
	return nil
}

func (u *User) CheckPassword(password string) bool {
	result, err := utils.ComparePasswordAndHash(password, *u.PasswordHash)
	if err != nil {
		logger.Error("Error during password hash compare!", err)
		return false
	}

	return result
}
