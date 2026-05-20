package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"

)

type UserRole string

const (
	RoleRegular        UserRole = "regular"
	RoleForeman        UserRole = "foreman"
	RoleProjectManager UserRole = "project_manager"
	RoleClerk          UserRole = "clerk"
	RoleEngineer       UserRole = "engineer"
	RoleDirector       UserRole = "director"
)

func (ct *UserRole) Scan(value any) error {
	*ct = UserRole(value.([]byte))
	return nil
}

func (ct UserRole) Value() (driver.Value, error) {
	return string(ct), nil
}

type User struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`

	Email        string   `gorm:"unique" json:"email"`
	PasswordHash string   `gorm:"column:password"`
	Role         UserRole `gorm:"type:user_role"`

	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt int64     `gorm:"autoUpdateTime:milli" json:"updated_at"`
}
