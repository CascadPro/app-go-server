package models

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;not null" json:"id"`

	Title  string   `gorm:"type:varchar(255);not null" json:"title"`
	Origin []string `gorm:"not null" json:"origin"`

	TechnicalDoc     uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_technical_doc_request" json:"technical_doc,omitempty"`
	ProjectDoc       uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_project_doc_request" json:"project_doc,omitempty"`
	SpecificationDoc uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_specification_doc_request" json:"specification_doc,omitempty"`

	WorkTypes []string `json:"work_types,omitempty"`
	Geography []string `json:"geography,omitempty"`

	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime:unix" json:"updated_at,omitempty"`
}
