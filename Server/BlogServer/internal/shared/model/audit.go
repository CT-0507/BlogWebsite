package model

import (
	"time"
)

type Audit struct {
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy *string    `json:"createdBy,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
	DeletedAt *time.Time `json:"deletedAt"`
	DeletedBy *string    `json:"deletedBy,omitempty"`
}
