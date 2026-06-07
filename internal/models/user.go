package models

import (
	"cascade/pkg/logger"
	"cascade/pkg/utils"
	"fmt"
	"path"
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

	Name     string  `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	Surname  string  `gorm:"type:varchar(255);not null" json:"surname,omitempty"`
	LastName *string `gorm:"type:varchar(255)" json:"last_name,omitempty"`

	Avatar *string `gorm:"type:uuid;uniqueIndex:idx_avatar_user" json:"avatar_path,omitempty"`

	Token []Token `gorm:"foreignKey:UserID" json:"tokens,omitempty"`

	LastActiveAt time.Time `gorm:"default:current_timestamp" json:"last_active_at,omitempty"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt    int64     `gorm:"autoUpdateTime:unix" json:"updated_at,omitempty"`
}

func (u *User) HashPassword() error {
	passwordHash, err := utils.GenerateHash(*u.PasswordHash)
	if err != nil {
		return err
	}

	u.PasswordHash = &passwordHash
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

func (u *User) SetAvatarPath(cfg utils.Config) {
	if u.Avatar == nil {
		return
	}

	avatarPath := path.Join(cfg.ApplicationUrl, utils.TagAvatars, *u.Avatar)
	u.Avatar = &avatarPath
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s %s", u.Name, u.Surname, *u.LastName)
}
