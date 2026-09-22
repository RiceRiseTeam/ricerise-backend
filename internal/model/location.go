package model

import (
	"time"

	"gorm.io/gorm"
)

type LocationModel struct {
	ID uint64 `gorm:"primary_key"`

	Name          string         `gorm:"not null;size:128"`
	Address       string         `gorm:"not null;size:256"`
	Longitude     float64        `gorm:"type:decimal(10,6)"`
	Latitude      float64        `gorm:"type:decimal(10,6)"`
	Reviewed      *bool          `gorm:"not null;default:false"`
	UserId        uint64         `gorm:"index;not null"`
	User          UserModel      `gorm:"foreignkey:UserId"`
	Description   string         `gorm:"type:text;not null"`
	AverageRating float32        `gorm:"type:decimal(3,2);not null;default:0"`
	Comments      []CommentModel `gorm:"foreignkey:LocationId"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
