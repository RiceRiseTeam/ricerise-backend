package model

import (
	"time"

	"gorm.io/gorm"
)

type LocationModel struct {
	ID uint64 `gorm:"primary_key"`

	Name        string  `gorm:"not null;size:128"`
	Longitude   float64 `gorm:"type:decimal(10,6)"`
	Latitude    float64 `gorm:"type:decimal(10,6)"`
	Reviewed    *bool   `gorm:"not null;default:false"`
	Description string  `gorm:"type:text;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
