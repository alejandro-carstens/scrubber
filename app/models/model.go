package models

import "time"

// Model represents the base model
type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key; type: bigint unsigned AUTO_INCREMENT;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
