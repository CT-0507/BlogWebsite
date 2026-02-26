package model

import (
	"time"
)

type Audit struct {
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *string    `json:"created_by,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy *string    `json:"deleted_by,omitempty"`
}
